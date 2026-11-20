package payments

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/publishing"
	"github.com/Includeoyasi/faultline-lab/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twmb/franz-go/pkg/kgo"
)

func Execute(ctx context.Context, args []string, lookup func(string) (string, bool), output io.Writer) int {
	if len(args) != 1 || args[0] != "process" {
		return service.Execute(ctx, "payments", args, lookup, output)
	}
	logger := slog.New(slog.NewJSONHandler(output, nil)).With("service", "payments")
	if err := run(ctx, lookup, output, logger); err != nil {
		logger.ErrorContext(ctx, "payment processing stopped", "error", err)
		return 1
	}
	return 0
}

func run(ctx context.Context, lookup func(string) (string, bool), output io.Writer, logger *slog.Logger) error {
	get := func(key string) string { v, _ := lookup(key); return v }
	if _, err := service.LoadConfig(lookup); err != nil {
		return err
	}
	dataset, seed := get("FL_PAYMENTS_DATASET_ID"), get("FL_PAYMENTS_SEED")
	bps, e1 := strconv.Atoi(get("FL_PAYMENTS_DECLINE_BPS"))
	delay, e2 := time.ParseDuration(get("FL_PAYMENTS_MAX_DELAY"))
	broker := get("FL_OUTBOX_BROKER")
	if !identifier.MatchString(dataset) || seed == "" || e1 != nil || bps < 0 || bps > 10000 || e2 != nil || delay < 0 || delay > time.Hour || broker == "" {
		return errors.New("invalid payment simulator or broker configuration")
	}
	host := get("FL_PAYMENTS_DB_ADDR")
	if h, port, err := net.SplitHostPort(host); err != nil || h == "" || port == "" {
		return errors.New("invalid Payments database address")
	}
	secret, err := os.ReadFile(get("FL_PAYMENTS_PASSWORD_FILE"))
	if err != nil {
		return errors.New("payment database credential unavailable")
	}
	addr := url.URL{Scheme: "postgres", Host: host, Path: "payments_db", User: url.UserPassword("payments_app", strings.TrimSpace(string(secret))), RawQuery: "sslmode=disable&require_auth=scram-sha-256&connect_timeout=5&statement_timeout=5000&pool_max_conns=4"}
	pool, err := pgxpool.New(ctx, addr.String())
	if err != nil {
		return errors.New("invalid payment database configuration")
	}
	defer pool.Close()
	p := processor{pool: pool, dataset: dataset, seed: seed, declineBPS: bps, maxDelay: delay}
	if err = p.bind(ctx); err != nil {
		return errors.New("cannot bind payment database to dataset")
	}
	consumer, err := newConsumer(broker, dataset, "orders.events.v1")
	if err != nil {
		return err
	}
	defer consumer.Close()
	workers, stop := context.WithCancel(ctx)
	defer stop()
	var wg sync.WaitGroup
	wg.Go(func() {
		publishing.Publisher{Pool: pool, Owner: "payments", Dataset: dataset, Broker: broker}.Run(workers, logger)
	})
	wg.Go(func() { p.work(workers, logger) })
	var receiveErr error
	wg.Go(func() { receiveErr = p.receive(workers, consumer); stop() })
	err = service.Run(workers, "payments", lookup, output)
	stop()
	wg.Wait()
	if ctx.Err() == nil && receiveErr != nil && !errors.Is(receiveErr, context.Canceled) {
		return receiveErr
	}
	return err
}

func newConsumer(broker, dataset, topic string) (*kgo.Client, error) {
	return kgo.NewClient(kgo.SeedBrokers(broker), kgo.ConsumeTopics(topic), kgo.ConsumerGroup("payments-v1-"+dataset), kgo.FetchIsolationLevel(kgo.ReadCommitted()), kgo.DisableAutoCommit(), kgo.BlockRebalanceOnPoll())
}

func (p processor) receive(ctx context.Context, cl *kgo.Client) error {
	defer cl.AllowRebalance()
	for ctx.Err() == nil {
		fetch := cl.PollRecords(ctx, 1)
		if errs := fetch.Errors(); len(errs) > 0 {
			return errs[0].Err
		}
		for _, r := range fetch.Records() {
			if err := p.save(ctx, r); err != nil {
				return err
			}
			if err := cl.CommitRecords(ctx, r); err != nil {
				return err
			}
		}
		cl.AllowRebalance()
	}
	return ctx.Err()
}

func (p processor) work(ctx context.Context, logger *slog.Logger) {
	for ctx.Err() == nil {
		err := p.finish(ctx)
		if err == nil {
			continue
		}
		if !errors.Is(err, pgx.ErrNoRows) && ctx.Err() == nil {
			logger.WarnContext(ctx, "payment completion delayed", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(100 * time.Millisecond):
		}
	}
}
