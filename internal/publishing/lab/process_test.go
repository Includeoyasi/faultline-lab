package lab

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/publishing"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
)

// Pause the child after it has read the pending fact under its own SQL lease.
type factReadBarrier struct{}
type pendingRead struct{}

func (factReadBarrier) TraceQueryStart(ctx context.Context, _ *pgx.Conn, q pgx.TraceQueryStartData) context.Context {
	if strings.HasPrefix(q.SQL, "SELECT event_id,order_id,aggregate_version,stream_sequence") {
		return context.WithValue(ctx, pendingRead{}, true)
	}
	return ctx
}
func (factReadBarrier) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, q pgx.TraceQueryEndData) {
	if ctx.Value(pendingRead{}) == true && q.Err == nil {
		fmt.Println("fact acquired")
		var signal [1]byte
		if _, err := io.ReadFull(os.Stdin, signal[:]); err != nil {
			panic(err)
		}
	}
}

func TestPausedPublisherProcess(t *testing.T) {
	if os.Getenv("FL_PAUSED_PUBLISHER_CHILD") != "1" {
		t.Skip("subprocess of the publisher takeover test")
	}
	ctx := t.Context()
	cfg, err := pgxpool.ParseConfig(os.Getenv("FL_PUBLISHING_TEST_DB"))
	if err != nil {
		t.Fatal(err)
	}
	cfg.ConnConfig.Tracer = factReadBarrier{}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	p := publishing.Publisher{Pool: pool, Dataset: "fixture", Owner: "inventory", Broker: os.Getenv("FL_PUBLISHING_TEST_BROKER")}
	err = p.Once(ctx)
	if !errors.Is(err, kerr.ProducerFenced) && !errors.Is(err, kerr.InvalidProducerEpoch) && !errors.Is(err, kerr.InvalidProducerIDMapping) {
		t.Fatalf("expected broker fencing, got %v", err)
	}
}

func verifyProcessTakeover(t *testing.T, ctx context.Context, pool, reader *pgxpool.Pool, p publishing.Publisher) {
	t.Helper()
	execSQL := func(q string) {
		t.Helper()
		if _, err := pool.Exec(ctx, q); err != nil {
			t.Fatal(err)
		}
	}
	execSQL("INSERT INTO inventory.outbox(event_id,order_id,aggregate_version,event_type,payload) VALUES ('three','o',7,'InventoryReserved','{}'),('four','o',8,'InventoryReleased','{}')")
	var initial int64
	if err := pool.QueryRow(ctx, "SELECT generation FROM inventory.publisher_lease").Scan(&initial); err != nil {
		t.Fatal(err)
	}
	childCtx, stop := context.WithTimeout(ctx, 18*time.Second)
	defer stop()
	cmd := exec.CommandContext(childCtx, os.Args[0], "-test.run=^TestPausedPublisherProcess$")
	cmd.Env = append(os.Environ(), "FL_PAUSED_PUBLISHER_CHILD=1")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}
	waited := false
	defer func() {
		stop()
		if !waited {
			_ = cmd.Wait()
		}
	}()
	out := bufio.NewReader(stdout)
	line, err := out.ReadString('\n')
	if err != nil || strings.TrimSpace(line) != "fact acquired" {
		t.Fatalf("pause barrier: %q %v", line, err)
	}
	// Advance only this disposable lease's deadline, preserving its generation.
	tag, err := pool.Exec(ctx, "UPDATE inventory.publisher_lease SET lease_until=clock_timestamp() WHERE generation=$1 AND lease_until>clock_timestamp()", initial+1)
	if err != nil || tag.RowsAffected() != 1 {
		t.Fatalf("old lease missing: %v", err)
	}
	for range 2 {
		if err = p.Once(ctx); err != nil {
			t.Fatal(err)
		}
	}
	if _, err = stdin.Write([]byte{1}); err != nil {
		t.Fatal(err)
	}
	_ = stdin.Close()
	tail, err := io.ReadAll(out)
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Wait()
	waited = true
	if err != nil {
		t.Fatalf("stale publisher: %s (%v)", tail, err)
	}
	var generation int64
	if err = pool.QueryRow(ctx, "SELECT generation FROM inventory.publisher_lease").Scan(&generation); err != nil || generation != initial+3 {
		t.Fatalf("unexpected generation %d: %v", generation, err)
	}
	// A committed marker on the same key bounds the log after the child exits.
	execSQL("INSERT INTO inventory.outbox(event_id,order_id,aggregate_version,event_type,payload) VALUES ('five','o',9,'InventoryCommitted','{}')")
	if err = p.Once(ctx); err != nil {
		t.Fatal(err)
	}
	cl, err := kgo.NewClient(kgo.SeedBrokers(p.Broker), kgo.ConsumeTopics("inventory.events.v1"), kgo.FetchIsolationLevel(kgo.ReadCommitted()), kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()))
	if err != nil {
		t.Fatal(err)
	}
	defer cl.Close()
	for _, want := range []int64{1, 1, 2, 3, 4, 5} {
		fetch := cl.PollRecords(ctx, 1)
		if es := fetch.Errors(); len(es) > 0 {
			t.Fatal(es[0].Err)
		}
		records := fetch.Records()
		if len(records) != 1 {
			t.Fatal("missing committed record")
		}
		r := records[0]
		var e publishing.Event
		if err = json.Unmarshal(r.Value, &e); err != nil || e.StreamSequence != want || e.OrderID != "o" {
			t.Fatalf("expected sequence %d, got %s", want, r.Value)
		}
		if err = Save(ctx, reader, p.Dataset, r); err != nil {
			t.Fatal(err)
		}
		if want >= 3 {
			if err = Apply(ctx, reader, p.Dataset, p.Owner, "o"); err != nil {
				t.Fatal(err)
			}
		}
	}
	var cursor, effects int
	var state string
	if err = reader.QueryRow(ctx, "SELECT cursor,effects,state FROM publishing_receiver.projection WHERE dataset='fixture' AND owner='inventory' AND order_id='o'").Scan(&cursor, &effects, &state); err != nil || cursor != 5 || effects != 5 || state != "InventoryCommitted" {
		t.Fatalf("final projection: %d %d %s (%v)", cursor, effects, state, err)
	}
}
