// Package store owns Orders data and short transactions, never peer RPCs.
package store

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/orders/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	bindDatasetSQL = "INSERT INTO orders.dataset(dataset_id) VALUES($1) ON CONFLICT(singleton) DO NOTHING"
	readDatasetSQL = "SELECT dataset_id FROM orders.dataset WHERE singleton"
	readOrderSQL   = `SELECT o.order_id,o.request_key,o.reserve_operation_id,o.state,o.version,o.items,o.reservation,o.rejection,j.blocked_reason
 FROM orders.orders o JOIN orders.reserve_jobs j USING(order_id)`
)

type Store struct {
	pool    *pgxpool.Pool
	dataset string
	timeout time.Duration
}

func NewStore(ctx context.Context, pool *pgxpool.Pool, dataset string, timeout time.Duration) (*Store, error) {
	if !model.ValidIdentifier(dataset) || timeout < time.Millisecond || timeout > time.Minute {
		return nil, model.ErrInvalid
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, storageError(err)
	}
	defer rollback(tx)
	if _, err = tx.Exec(ctx, bindDatasetSQL, dataset); err != nil {
		return nil, storageError(err)
	}
	var actual string
	if err = tx.QueryRow(ctx, readDatasetSQL).Scan(&actual); err != nil {
		return nil, storageError(err)
	}
	if actual != dataset {
		return nil, model.ErrDataset
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, storageError(err)
	}
	return &Store{pool: pool, dataset: dataset, timeout: timeout}, nil
}
func rollback(tx pgx.Tx) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = tx.Rollback(ctx)
}
func storageError(err error) error {
	for _, known := range []error{context.Canceled, context.DeadlineExceeded, model.ErrNotFound, model.ErrStale} {
		if errors.Is(err, known) {
			return known
		}
	}
	return errors.New("orders storage operation failed; outcome may be unknown")
}
func (s *Store) read(row pgx.Row) (model.Order, error) {
	o := model.Order{DatasetID: s.dataset}
	var items, snapshot []byte
	err := row.Scan(&o.OrderID, &o.RequestKey, &o.OperationID, &o.State, &o.Version, &items, &snapshot, &o.Rejection, &o.BlockedReason)
	if errors.Is(err, pgx.ErrNoRows) {
		return o, model.ErrNotFound
	}
	if err != nil {
		return o, storageError(err)
	}
	if err = json.Unmarshal(items, &o.Items); err != nil {
		return o, storageError(err)
	}
	if len(snapshot) > 0 {
		if err = json.Unmarshal(snapshot, &o.Reservation); err != nil {
			return o, storageError(err)
		}
	}
	return o, nil
}
func (s *Store) validateLookup(dataset, key string) error {
	if !model.ValidIdentifier(dataset) || !model.ValidIdentifier(key) {
		return model.ErrInvalid
	}
	if dataset != s.dataset {
		return model.ErrDataset
	}
	return nil
}
