package store

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Includeoyasi/faultline-lab/internal/orders/model"
	"github.com/Includeoyasi/faultline-lab/internal/publishing"
	"github.com/jackc/pgx/v5"
	"github.com/twmb/franz-go/pkg/kgo"
)

// ReceiveEvent implements the transactional inbox. A Kafka offset is committed
// only after the complete event and its contiguous source sequence are durable.
func (s *Store) ReceiveEvent(ctx context.Context, record *kgo.Record) error {
	var event publishing.Event
	if err := json.Unmarshal(record.Value, &event); err != nil {
		return errors.New("invalid event JSON")
	}
	if event.Dataset != s.dataset || !model.ValidIdentifier(event.OrderID) ||
		event.EventID == "" || len(event.EventID) > 256 || string(record.Key) != event.OrderID ||
		event.SchemaVersion != 1 || event.AggregateVersion < 1 || event.StreamSequence < 1 ||
		event.OccurredAt.IsZero() || record.Topic != event.Owner+".events.v1" ||
		(event.Owner != "payments" && event.Owner != "inventory") {
		return errors.New("invalid event envelope")
	}
	switch event.Type {
	case "PaymentSucceeded", "PaymentDeclined":
		if event.Owner != "payments" {
			return errors.New("event owner and type disagree")
		}
	case "InventoryReserved", "InventoryCommitted", "InventoryReleased":
		if event.Owner != "inventory" {
			return errors.New("event owner and type disagree")
		}
	default:
		return errors.New("unsupported event type")
	}
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO orders.inbox_cursor(owner,order_id,sequence)
 VALUES($1,$2,0) ON CONFLICT DO NOTHING`, event.Owner, event.OrderID)
		if err != nil {
			return storageError(err)
		}
		var sequence int64
		if err = tx.QueryRow(ctx, "SELECT sequence FROM orders.inbox_cursor WHERE owner=$1 AND order_id=$2 FOR UPDATE",
			event.Owner, event.OrderID).Scan(&sequence); err != nil {
			return storageError(err)
		}
		var saved []byte
		err = tx.QueryRow(ctx, "SELECT body FROM orders.inbox WHERE owner=$1 AND event_id=$2",
			event.Owner, event.EventID).Scan(&saved)
		if err == nil {
			if !bytes.Equal(saved, record.Value) {
				return errors.New("event identity has conflicting content")
			}
			return nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return storageError(err)
		}
		if event.StreamSequence != sequence+1 {
			return fmt.Errorf("event sequence gap or conflict: expected %d, received %d", sequence+1, event.StreamSequence)
		}
		_, err = tx.Exec(ctx, `INSERT INTO orders.inbox(owner,event_id,order_id,stream_sequence,body)
 VALUES($1,$2,$3,$4,$5)`, event.Owner, event.EventID, event.OrderID, event.StreamSequence, record.Value)
		if err != nil {
			return storageError(err)
		}
		_, err = tx.Exec(ctx, "UPDATE orders.inbox_cursor SET sequence=$3 WHERE owner=$1 AND order_id=$2",
			event.Owner, event.OrderID, event.StreamSequence)
		if err != nil {
			return storageError(err)
		}
		return nil
	})
}

// deferInbox preserves an event until its prerequisite arrives or the conflict
// can be inspected. It never changes a business state on an unverified fact.
func deferInbox(ctx context.Context, tx pgx.Tx, owner, eventID, reason string) error {
	_, err := tx.Exec(ctx, `UPDATE orders.inbox SET available_at=clock_timestamp()+interval '1 second',
 last_error=$3 WHERE owner=$1 AND event_id=$2`, owner, eventID, reason)
	return err
}

// ApplyEvent applies the next saved fact for one producer stream in a short
// transaction. Events from Payments and Inventory may arrive in either order.
func (s *Store) ApplyEvent(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	worked := false
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var owner, eventID string
		var body []byte
		err := tx.QueryRow(ctx, `SELECT i.owner,i.event_id,i.body FROM orders.inbox i
 WHERE NOT i.applied AND i.available_at<=clock_timestamp()
 AND NOT EXISTS (SELECT FROM orders.inbox p WHERE p.owner=i.owner AND p.order_id=i.order_id
                 AND NOT p.applied AND p.stream_sequence<i.stream_sequence)
 ORDER BY i.available_at,i.owner,i.order_id,i.stream_sequence
 LIMIT 1 FOR UPDATE OF i SKIP LOCKED`).Scan(&owner, &eventID, &body)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		worked = true
		var event publishing.Event
		if err = json.Unmarshal(body, &event); err != nil {
			return deferInbox(ctx, tx, owner, eventID, "INVALID_SAVED_EVENT")
		}
		order, err := s.read(tx.QueryRow(ctx, lockOrderSQL, event.OrderID))
		if errors.Is(err, model.ErrNotFound) {
			return deferInbox(ctx, tx, owner, eventID, "ORDER_NOT_FOUND")
		}
		if err != nil {
			return err
		}
		var reason string
		switch owner {
		case "payments":
			reason, err = s.applyPayment(ctx, tx, order, event)
		case "inventory":
			reason, err = s.applyInventory(ctx, tx, order, event)
		}
		if err != nil {
			return err
		}
		if reason != "" {
			return deferInbox(ctx, tx, owner, eventID, reason)
		}
		_, err = tx.Exec(ctx, "UPDATE orders.inbox SET applied=true,last_error='' WHERE owner=$1 AND event_id=$2",
			owner, eventID)
		return err
	})
	if err != nil {
		return worked, storageError(err)
	}
	return worked, nil
}
