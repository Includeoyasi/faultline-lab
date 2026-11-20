package payments

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/publishing"
	"github.com/Includeoyasi/faultline-lab/migrations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twmb/franz-go/pkg/kgo"
)

// The first process exits with accepted offsets and no worker running.
func TestPaymentIntakeChild(t *testing.T) {
	if os.Getenv("FL_PAYMENT_INTAKE_CHILD") != "1" {
		t.Skip("started by recovery test")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 45*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, os.Getenv("FL_PUBLISHING_TEST_DB"))
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	p := processor{pool: pool, dataset: "payment-test", seed: "example", declineBPS: 5000, maxDelay: time.Hour}
	if err = p.bind(ctx); err != nil {
		t.Fatal(err)
	}
	cl, err := newConsumer(os.Getenv("FL_PUBLISHING_TEST_BROKER"), p.dataset, "payment-test-input")
	if err != nil {
		t.Fatal(err)
	}
	defer cl.Close()
	if err = p.receive(ctx, cl); err == nil || err.Error() != "invalid Orders event JSON" {
		t.Fatalf("unexpected intake exit: %v", err)
	}
}

func TestPaymentRecovery(t *testing.T) {
	dsn, broker := os.Getenv("FL_PUBLISHING_TEST_DB"), os.Getenv("FL_PUBLISHING_TEST_BROKER")
	if dsn == "" || broker == "" {
		t.Skip("requires isolated Kafka and PostgreSQL")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 80*time.Second)
	defer cancel()
	check := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	connect := func() *pgxpool.Pool { pool, err := pgxpool.New(ctx, dsn); check(err); return pool }
	first := connect()
	defer first.Close()
	paths, err := fs.Glob(migrations.Files, "payments/*.sql")
	check(err)
	for _, path := range paths {
		b, err := migrations.Files.ReadFile(path)
		check(err)
		_, err = first.Exec(ctx, string(b))
		check(err)
	}
	p := processor{pool: first, dataset: "payment-test", seed: "example", declineBPS: 5000, maxDelay: time.Hour}
	check(p.bind(ctx))
	other := p
	other.dataset = "other"
	if other.bind(ctx) == nil {
		t.Fatal("bound to another dataset")
	}
	const topic = "payment-test-input"
	producer, err := kgo.NewClient(kgo.SeedBrokers(broker), kgo.ProducerBatchCompression(kgo.NoCompression()), kgo.TransactionalID("orders-payment-test-input"))
	check(err)
	defer producer.Close()
	expected := map[string]string{}
	var records []*kgo.Record
	seen := map[bool]bool{}
	for i := 0; len(records) < 2 && i < 1000; i++ {
		order := fmt.Sprintf("order-%d", i)
		approved, delay := p.simulate(order)
		if seen[approved] || delay < time.Minute {
			continue
		}
		seen[approved] = true
		kind := "PaymentDeclined"
		if approved {
			kind = "PaymentSucceeded"
		}
		expected[order] = kind
		payload := fmt.Sprintf(`{"dataset_id":"payment-test","order_id":%q,"version":3,"state":"PAYMENT_PENDING","reservation":{"total_minor":1250,"currency":"RUB"}}`, order)
		e := publishing.Event{Owner: "orders", Dataset: p.dataset, OrderID: order, EventID: order + ":awaiting-payment", Type: "OrderAwaitingPayment", StreamSequence: 1, AggregateVersion: 3, SchemaVersion: 1, OccurredAt: time.Now().UTC(), Payload: json.RawMessage(payload)}
		value, err := json.Marshal(e)
		check(err)
		// A new consumer must recover a backlog older than the client's
		// default reset window.
		records = append(records, &kgo.Record{Topic: topic, Key: []byte(order), Value: value, Timestamp: time.Now().Add(-2 * time.Hour)})
	}
	if len(records) != 2 {
		t.Fatal("missing simulator outcome")
	}
	// Mirror the real Orders relay: initialize the transactional producer
	// before publishing the batch that a read_committed consumer can see.
	_, _, err = producer.ProducerID(ctx)
	check(err)
	check(producer.BeginTransaction())
	for _, r := range records {
		check(producer.ProduceSync(ctx, r).FirstErr())
	}
	check(producer.ProduceSync(ctx, &kgo.Record{Topic: topic, Key: records[0].Key, Value: records[0].Value, Timestamp: records[0].Timestamp}).FirstErr())
	// A fresh group consumes until a sentinel on the same partition. Its offset
	// remains unacknowledged; both preceding jobs must already be durable.
	sentinel := &kgo.Record{Topic: topic, Value: []byte("invalid")}
	check(producer.ProduceSync(ctx, sentinel).FirstErr())
	check(producer.EndTransaction(ctx, kgo.TryCommit))
	consume := func(p processor) {
		cl, err := newConsumer(broker, p.dataset, topic)
		check(err)
		defer cl.Close()
		err = p.receive(ctx, cl)
		if err == nil || err.Error() != "invalid Orders event JSON" {
			t.Fatalf("expected unacknowledged invalid input, got %v", err)
		}
	}
	child := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestPaymentIntakeChild$")
	child.Env = append(os.Environ(), "FL_PAYMENT_INTAKE_CHILD=1")
	if out, err := child.CombinedOutput(); err != nil {
		t.Fatalf("intake process: %v: %s", err, out)
	}
	if err = p.finish(ctx); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("simulated delay: %v", err)
	}
	first.Close()
	p.pool = connect()
	defer p.pool.Close()
	check(p.bind(ctx))
	consume(p)
	for _, r := range records {
		check(p.save(ctx, r))
	}
	var changed publishing.Event
	check(json.Unmarshal(records[0].Value, &changed))
	changed.EventID += ":conflict"
	conflicting, err := json.Marshal(changed)
	check(err)
	if err = p.save(ctx, &kgo.Record{Key: records[0].Key, Value: conflicting}); err == nil {
		t.Fatal("conflicting request accepted")
	}
	// Advance only the persisted test deadline, without sleeping for business time.
	_, err = p.pool.Exec(ctx, "UPDATE payments.jobs SET ready_at=clock_timestamp()")
	check(err)
	for range records {
		check(p.finish(ctx))
	}
	for _, r := range records {
		check(p.save(ctx, r))
	}
	if err = p.finish(ctx); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("extra job: %v", err)
	}
	pub := publishing.Publisher{Pool: p.pool, Owner: "payments", Dataset: p.dataset, Broker: broker}
	for range records {
		check(pub.Once(ctx))
	}
	reader, err := kgo.NewClient(kgo.SeedBrokers(broker), kgo.ConsumeTopics("payments.events.v1"), kgo.FetchIsolationLevel(kgo.ReadCommitted()))
	check(err)
	defer reader.Close()
	for len(expected) > 0 {
		f := reader.PollRecords(ctx, 1)
		if errs := f.Errors(); len(errs) > 0 {
			t.Fatal(errs[0].Err)
		}
		for _, r := range f.Records() {
			var e publishing.Event
			check(json.Unmarshal(r.Value, &e))
			if e.Dataset != p.dataset {
				continue
			}
			if expected[e.OrderID] != e.Type || e.Type == "" || string(r.Key) != e.OrderID || e.Dataset != p.dataset || e.Owner != "payments" || e.StreamSequence != 1 {
				t.Fatalf("wrong result: %+v", e)
			}
			var count int
			var net int64
			check(p.pool.QueryRow(ctx, "SELECT count(*),COALESCE(sum(amount_minor),0)::bigint FROM payments.payment_ledger WHERE order_id=$1 AND currency='RUB'", e.OrderID).Scan(&count, &net))
			if (e.Type == "PaymentSucceeded" && (count != 1 || net != 1250)) || (e.Type == "PaymentDeclined" && (count != 0 || net != 0)) {
				t.Fatalf("%s: %d ledger rows, %d minor units", e.Type, count, net)
			}
			delete(expected, e.OrderID)
		}
	}
}
