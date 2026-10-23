package store

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/orders/model"
	"github.com/jackc/pgx/v5"
)

const (
	lockOwnedJobSQL = `SELECT order_version FROM orders.reserve_jobs WHERE order_id=$1 AND generation=$2
 AND NOT done AND blocked_reason='' AND lease_until>clock_timestamp() FOR UPDATE`
	resolveOrderSQL = `UPDATE orders.orders SET state=$3,version=version+1,reservation=$4,rejection=$5,updated_at=clock_timestamp()
 WHERE order_id=$1 AND version=$2 AND state='RESERVING' RETURNING version`
	insertEventSQL = `INSERT INTO orders.outbox(event_id,order_id,aggregate_version,event_type,payload) VALUES($1,$2,$3,'OrderAwaitingPayment',$4)`
	resolveJobSQL  = `UPDATE orders.reserve_jobs SET done=$3,blocked_reason=$4,last_error=$5,
 available_at=clock_timestamp()+make_interval(secs=>$6),lease_until=NULL
 WHERE order_id=$1 AND generation=$2 AND lease_until>clock_timestamp()`
)

// Resolve atomically commits a verified result or reschedules the same operation.
// Job-first row locking matches Claim. The final lease CAS also covers lock waits.
func (s *Store) Resolve(ctx context.Context, j model.Job, r model.Resolution, delay time.Duration) error {
	count := 0
	for _, set := range []bool{r.Reservation != nil, r.Rejection != "", r.BlockedReason != "", r.RetryReason != ""} {
		if set {
			count++
		}
	}
	if count != 1 || delay < 0 || delay > time.Minute || !model.ValidIdentifier(j.Order.OrderID) || j.Generation < 1 {
		return model.ErrInvalid
	}
	if j.Order.DatasetID != s.dataset {
		return model.ErrDataset
	}
	if r.Rejection != "" && r.Rejection != "INSUFFICIENT_STOCK" && r.Rejection != "UNKNOWN_SKU" {
		return model.ErrInvalid
	}
	if len(r.BlockedReason) > 128 || len(r.RetryReason) > 128 {
		return model.ErrInvalid
	}
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return storageError(err)
	}
	defer rollback(tx)
	var expected int64
	if err = tx.QueryRow(ctx, lockOwnedJobSQL, j.Order.OrderID, j.Generation).Scan(&expected); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrStale
		}
		return storageError(err)
	}
	current, err := s.read(tx.QueryRow(ctx, lockOrderSQL, j.Order.OrderID))
	if err != nil {
		return err
	}
	if expected != j.OrderVersion || current.Version != j.OrderVersion || current.State != model.Reserving {
		return model.ErrStale
	}
	done := r.Reservation != nil || r.Rejection != ""
	if done {
		current.Version++
		current.Rejection = r.Rejection
		current.Reservation = r.Reservation
		current.State = model.Rejected
		var snapshot []byte
		if r.Reservation != nil {
			current.State = model.PaymentPending
			snapshot, err = json.Marshal(r.Reservation)
			if err != nil {
				return err
			}
		}
		if err = tx.QueryRow(ctx, resolveOrderSQL, current.OrderID, j.OrderVersion, current.State, snapshot, r.Rejection).Scan(&current.Version); err != nil {
			return storageError(err)
		}
		if r.Reservation != nil {
			payload, e := json.Marshal(current)
			if e != nil {
				return e
			}
			if _, err = tx.Exec(ctx, insertEventSQL, current.OperationID+":awaiting-payment", current.OrderID, current.Version, payload); err != nil {
				return storageError(err)
			}
		}
	}
	tag, err := tx.Exec(ctx, resolveJobSQL, current.OrderID, j.Generation, done, r.BlockedReason, r.RetryReason, delay.Seconds())
	if err != nil {
		return storageError(err)
	}
	if tag.RowsAffected() != 1 {
		return model.ErrStale
	}
	if err = tx.Commit(ctx); err != nil {
		return storageError(err)
	}
	return nil
}
