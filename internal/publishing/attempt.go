// Package publishing delivers durable outbox facts using a single-use producer.
package publishing

import (
	"context"
	"errors"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

// attempt owns the client until the caller returns from the SQL mark.
// Cancellation also terminates in-flight network operations.
func attempt(ctx context.Context, brokers []string, lane string, hooks []kgo.Hook,
	claim func(context.Context) (*kgo.Record, error), mark func(context.Context) error) error {
	cl, err := kgo.NewClient(kgo.WithContext(ctx), kgo.SeedBrokers(brokers...), kgo.TransactionalID(lane),
		kgo.TransactionTimeout(15*time.Second), kgo.WithHooks(hooks...))
	if err != nil {
		return err
	}
	defer cl.Close()
	id, _, err := cl.ProducerID(ctx)
	if err != nil {
		return err
	}
	if id < 0 {
		return errors.New("broker did not initialize transactional producer")
	}
	record, err := claim(ctx)
	if err != nil {
		return err
	}
	if err = cl.BeginTransaction(); err != nil {
		return err
	}
	if err = cl.ProduceSync(ctx, record).FirstErr(); err != nil {
		return err
	}
	if err = cl.EndTransaction(ctx, kgo.TryCommit); err != nil {
		return err
	}
	return mark(ctx)
}
