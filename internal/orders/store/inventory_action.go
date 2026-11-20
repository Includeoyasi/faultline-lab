package store

import (
	"context"
	"errors"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/orders/model"
	"github.com/jackc/pgx/v5"
)

// ClaimInventoryAction uses a lease and generation so another worker can
// resume the same operation after a crash.
func (s *Store) ClaimInventoryAction(ctx context.Context, lease time.Duration) (*model.InventoryAction, error) {
	if lease < time.Millisecond || lease > time.Minute {
		return nil, model.ErrInvalid
	}
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	var result *model.InventoryAction
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var id, operation, action string
		err := tx.QueryRow(ctx, `SELECT order_id,operation_id,action FROM orders.inventory_actions
 WHERE NOT done AND blocked_reason='' AND available_at<=clock_timestamp()
 AND (lease_until IS NULL OR lease_until<=clock_timestamp())
 ORDER BY available_at,order_id LIMIT 1 FOR UPDATE SKIP LOCKED`).Scan(&id, &operation, &action)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		order, err := s.read(tx.QueryRow(ctx, readOrderSQL+" WHERE o.order_id=$1", id))
		if err != nil {
			return err
		}
		if (action == "commit" && order.State != model.Finalizing) ||
			(action == "release" && order.State != model.Cancelling) {
			return errors.New("inventory action has incompatible order state")
		}
		job := &model.InventoryAction{Order: order, OperationID: operation, Action: action}
		err = tx.QueryRow(ctx, `UPDATE orders.inventory_actions
 SET generation=generation+1,attempts=attempts+1,lease_until=clock_timestamp()+make_interval(secs=>$2)
 WHERE order_id=$1 RETURNING generation,attempts`, id, lease.Seconds()).Scan(&job.Generation, &job.Attempts)
		if err != nil {
			return err
		}
		result = job
		return nil
	})
	if err != nil {
		return nil, storageError(err)
	}
	return result, nil
}

// RecordInventoryAttempt does not finish the order. The corresponding
// Inventory event is the durable confirmation of the terminal state.
func (s *Store) RecordInventoryAttempt(ctx context.Context, job model.InventoryAction, blocked, reason string, delay time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	tag, err := s.pool.Exec(ctx, `UPDATE orders.inventory_actions SET lease_until=NULL,
 available_at=clock_timestamp()+make_interval(secs=>$4),blocked_reason=$5,last_error=$6
 WHERE order_id=$1 AND generation=$2 AND operation_id=$3 AND NOT done`,
		job.Order.OrderID, job.Generation, job.OperationID, delay.Seconds(), blocked, reason)
	if err != nil {
		return storageError(err)
	}
	if tag.RowsAffected() == 0 {
		return model.ErrStale
	}
	return nil
}
