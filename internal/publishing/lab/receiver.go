// Package lab verifies durable reception in an isolated test database.
package lab

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"

	"github.com/Includeoyasi/faultline-lab/internal/publishing"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twmb/franz-go/pkg/kgo"
)

//go:embed schema.sql
var Schema string

// Receive persists each partition prefix before acknowledging its offsets.
// The caller binds this client to an isolated group and read_committed.
func Receive(ctx context.Context, cl *kgo.Client, pool *pgxpool.Pool, dataset string) error {
	for ctx.Err() == nil {
		fetch := cl.PollRecords(ctx, 1)
		if errs := fetch.Errors(); len(errs) > 0 {
			return errs[0].Err
		}
		for _, r := range fetch.Records() {
			if err := Save(ctx, pool, dataset, r); err != nil {
				return err
			}
			if err := cl.CommitRecords(ctx, r); err != nil {
				return err
			}
		}
	}
	return ctx.Err()
}

// Save implements the Idempotent Consumer pattern by committing the inbox entry and job together.
func Save(ctx context.Context, pool *pgxpool.Pool, dataset string, r *kgo.Record) error {
	var e publishing.Event
	diagnostic := func(reason string) error {
		_, err := pool.Exec(ctx, "INSERT INTO publishing_receiver.diagnostic VALUES($1,$2,$3,$4) ON CONFLICT DO NOTHING", r.Topic, r.Partition, r.Offset, reason)
		if err != nil {
			return err
		}
		return errors.New(reason)
	}
	if json.Unmarshal(r.Value, &e) != nil || e.Dataset != dataset || (e.Owner != "orders" && e.Owner != "inventory") || r.Topic != e.Owner+".events.v1" || string(r.Key) != e.OrderID || e.OrderID == "" || e.EventID == "" || e.SchemaVersion != 1 || e.StreamSequence < 1 {
		return diagnostic("invalid event envelope")
	}
	conflict := errors.New("conflicting event identity or sequence")
	err := pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "INSERT INTO publishing_receiver.inbox VALUES($1,$2,$3,$4,$5,$6,$7) ON CONFLICT DO NOTHING", dataset, e.Owner, e.EventID, e.OrderID, e.StreamSequence, e.Type, r.Value)
		if err != nil {
			return err
		}
		var stored []byte
		err = tx.QueryRow(ctx, "SELECT body FROM publishing_receiver.inbox WHERE dataset=$1 AND owner=$2 AND event_id=$3", dataset, e.Owner, e.EventID).Scan(&stored)
		if errors.Is(err, pgx.ErrNoRows) || err == nil && !bytes.Equal(stored, r.Value) {
			return conflict
		}
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, "INSERT INTO publishing_receiver.job(dataset,owner,event_id) VALUES($1,$2,$3) ON CONFLICT DO NOTHING", dataset, e.Owner, e.EventID)
		return err
	})
	if errors.Is(err, conflict) {
		return diagnostic(conflict.Error())
	}
	return err
}

// Apply advances one stream in one SQL transaction. A missing predecessor
// or unknown type leaves the job pending for inspection and retry.
func Apply(ctx context.Context, pool *pgxpool.Pool, dataset, owner, order string) error {
	return pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		args := []any{dataset, owner, order}
		_, err := tx.Exec(ctx, "INSERT INTO publishing_receiver.projection(dataset,owner,order_id) VALUES($1,$2,$3) ON CONFLICT DO NOTHING", args...)
		if err != nil {
			return err
		}
		var cursor, seq int64
		var id, kind string
		err = tx.QueryRow(ctx, "SELECT cursor FROM publishing_receiver.projection WHERE dataset=$1 AND owner=$2 AND order_id=$3 FOR UPDATE", args...).Scan(&cursor)
		if err != nil {
			return err
		}
		err = tx.QueryRow(ctx, "SELECT i.event_id,i.kind,i.seq FROM publishing_receiver.inbox i JOIN publishing_receiver.job j USING(dataset,owner,event_id) WHERE dataset=$1 AND owner=$2 AND i.order_id=$3 AND NOT j.done ORDER BY seq LIMIT 1", args...).Scan(&id, &kind, &seq)
		if err != nil {
			return err
		}
		if seq != cursor+1 {
			return errors.New("missing predecessor")
		}
		relevant := owner == "inventory"
		if (relevant && kind != "InventoryReserved" && kind != "InventoryCommitted" && kind != "InventoryReleased") || (!relevant && kind != "OrderAwaitingPayment") {
			return errors.New("unknown event contract")
		}
		_, err = tx.Exec(ctx, "UPDATE publishing_receiver.projection SET cursor=$4,state=CASE WHEN $5 THEN $6 ELSE state END,effects=effects+CASE WHEN $5 THEN 1 ELSE 0 END WHERE dataset=$1 AND owner=$2 AND order_id=$3", dataset, owner, order, seq, relevant, kind)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, "UPDATE publishing_receiver.job SET done=true WHERE dataset=$1 AND owner=$2 AND event_id=$3", dataset, owner, id)
		return err
	})
}
