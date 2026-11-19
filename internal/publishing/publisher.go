package publishing

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Event struct {
	Dataset          string          `json:"dataset"`
	Owner            string          `json:"owner"`
	EventID          string          `json:"event_id"`
	OrderID          string          `json:"order_id"`
	AggregateVersion int64           `json:"aggregate_version"`
	StreamSequence   int64           `json:"stream_sequence"`
	Type             string          `json:"event_type"`
	Payload          json.RawMessage `json:"payload"`
	OccurredAt       time.Time       `json:"occurred_at"`
	SchemaVersion    int             `json:"schema_version"`
}

type Publisher struct {
	Pool                   *pgxpool.Pool
	Owner, Dataset, Broker string
}

// Once implements the Transactional Outbox relay and attempts one fact. Generation is observed before Kafka initialization.
// The mark and lease release share a SQL statement and lock the lease row.
func (p Publisher) Once(ctx context.Context) error {
	if p.Owner != "orders" && p.Owner != "inventory" && p.Owner != "payments" {
		return errors.New("unknown publisher schema")
	}
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	lease := p.Owner + ".publisher_lease"
	outbox := p.Owner + ".outbox"
	var g int64
	err := p.Pool.QueryRow(ctx, "SELECT generation FROM "+lease+" WHERE singleton AND lease_until <= clock_timestamp() AND EXISTS (SELECT FROM "+outbox+" WHERE published_at IS NULL)").Scan(&g)
	if err != nil {
		return err
	}
	event := Event{Owner: p.Owner, Dataset: p.Dataset, SchemaVersion: 1}
	return attempt(ctx, []string{p.Broker}, p.Dataset+"/"+p.Owner+"/0", nil, func(ctx context.Context) (*kgo.Record, error) {
		err := p.Pool.QueryRow(ctx, "UPDATE "+lease+" SET generation=generation+1, lease_until=clock_timestamp()+interval '25 seconds' WHERE singleton AND generation=$1 AND lease_until<=clock_timestamp() RETURNING generation", g).Scan(&g)
		if err != nil {
			return nil, err
		}
		err = p.Pool.QueryRow(ctx, "SELECT event_id,order_id,aggregate_version,stream_sequence,event_type,payload,created_at FROM "+outbox+" WHERE published_at IS NULL ORDER BY order_id,stream_sequence LIMIT 1").Scan(&event.EventID, &event.OrderID, &event.AggregateVersion, &event.StreamSequence, &event.Type, &event.Payload, &event.OccurredAt)
		if err != nil {
			return nil, err
		}
		event.OccurredAt = event.OccurredAt.UTC()
		value, err := json.Marshal(event)
		if err != nil {
			return nil, err
		}
		return &kgo.Record{Topic: p.Owner + ".events.v1", Key: []byte(event.OrderID), Value: value}, nil
	}, func(ctx context.Context) error {
		tag, err := p.Pool.Exec(ctx, "WITH released AS (UPDATE "+lease+" SET lease_until='-infinity' WHERE singleton AND generation=$1 AND lease_until>clock_timestamp() RETURNING singleton) UPDATE "+outbox+" SET published_at=clock_timestamp() FROM released WHERE event_id=$2 AND published_at IS NULL", g, event.EventID)
		if err != nil {
			return err
		}
		if tag.RowsAffected() != 1 {
			return errors.New("publication lease expired before mark")
		}
		return nil
	})
}

func (p Publisher) Run(ctx context.Context, logger *slog.Logger) {
	for ctx.Err() == nil {
		err := p.Once(ctx)
		if err == nil {
			continue
		}
		if !errors.Is(err, pgx.ErrNoRows) && ctx.Err() == nil {
			logger.WarnContext(ctx, "outbox publication delayed", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(250*time.Millisecond + rand.N(250*time.Millisecond)):
		}
	}
}
