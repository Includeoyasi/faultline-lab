package publishing

import (
	"context"
	"errors"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type initTrace struct{ count atomic.Int32 }

func (h *initTrace) OnBrokerWrite(_ kgo.BrokerMetadata, key int16, _ int, _, _ time.Duration, _ error) {
	if key == 22 {
		h.count.Add(1)
	}
}

// Run only against the disposable publishing test broker with its three-partition test topic.
func TestTakeover(t *testing.T) {
	broker := os.Getenv("FL_PUBLISHING_TEST_BROKER")
	if broker == "" {
		t.Skip("set FL_PUBLISHING_TEST_BROKER for isolated broker")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 45*time.Second)
	defer cancel()
	lane := "publisher-takeover-" + time.Now().UTC().Format("20060102150405.000000000")
	record := func(value string) *kgo.Record {
		return &kgo.Record{Topic: "orders.events.v1", Key: []byte(lane), Value: []byte(value)}
	}
	ready := make(chan struct{})
	resume := make(chan struct{})
	done := make(chan error, 1)
	var old, newTrace initTrace
	var marked atomic.Bool
	go func() {
		done <- attempt(ctx, []string{broker}, lane, []kgo.Hook{&old}, func(ctx context.Context) (*kgo.Record, error) {
			close(ready)
			select {
			case <-resume:
				return record("stale"), nil
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}, func(context.Context) error { marked.Store(true); return nil })
	}()
	select {
	case <-ready:
	case err := <-done:
		t.Fatalf("first init: %v", err)
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	err := attempt(ctx, []string{broker}, lane, []kgo.Hook{&newTrace}, func(context.Context) (*kgo.Record, error) { return record("replacement"), nil }, func(context.Context) error { return nil })
	close(resume)
	if err != nil {
		t.Fatal(err)
	}
	if err = <-done; err == nil || errors.Is(err, context.DeadlineExceeded) || marked.Load() {
		t.Fatalf("stale outcome: %v, marked %v", err, marked.Load())
	}
	if old.count.Load() != 1 || newTrace.count.Load() != 1 {
		t.Fatalf("unexpected InitProducerID counts: old %d replacement %d", old.count.Load(), newTrace.count.Load())
	}
}
