package store

import (
	"context"
	"errors"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/orders/model"
	"github.com/jackc/pgx/v5"
)

const (
	selectJobSQL = `SELECT j.order_id FROM orders.reserve_jobs j
 WHERE NOT j.done AND j.blocked_reason='' AND j.available_at<=clock_timestamp()
 AND (j.lease_until IS NULL OR j.lease_until<=clock_timestamp())
 ORDER BY j.available_at,j.order_id LIMIT 1 FOR UPDATE OF j SKIP LOCKED`
	lockOrderSQL  = readOrderSQL + " WHERE o.order_id=$1 FOR UPDATE OF o"
	startOrderSQL = "UPDATE orders.orders SET state='RESERVING',version=version+1,updated_at=clock_timestamp() WHERE order_id=$1 AND state='ACCEPTED' RETURNING version"
	claimJobSQL   = `UPDATE orders.reserve_jobs SET generation=generation+1,attempts=attempts+1,order_version=$2,
 lease_until=clock_timestamp()+make_interval(secs=>$3) WHERE order_id=$1 RETURNING generation,attempts`
)

func (s *Store) Claim(ctx context.Context, lease time.Duration) (*model.Job, error) {
	if lease < time.Millisecond || lease > time.Minute {
		return nil, model.ErrInvalid
	}
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, storageError(err)
	}
	defer rollback(tx)
	var id string
	err = tx.QueryRow(ctx, selectJobSQL).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, storageError(err)
	}
	o, err := s.read(tx.QueryRow(ctx, lockOrderSQL, id))
	if err != nil {
		return nil, err
	}
	if o.State == model.Accepted {
		if err = tx.QueryRow(ctx, startOrderSQL, id).Scan(&o.Version); err != nil {
			return nil, storageError(err)
		}
		o.State = model.Reserving
	}
	if o.State != model.Reserving {
		return nil, errors.New("reserve job has incompatible order phase")
	}
	job := &model.Job{Order: o, OrderVersion: o.Version}
	if err = tx.QueryRow(ctx, claimJobSQL, id, o.Version, lease.Seconds()).Scan(&job.Generation, &job.Attempts); err != nil {
		return nil, storageError(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, storageError(err)
	}
	return job, nil
}
