package store

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/Includeoyasi/faultline-lab/internal/orders/model"
	"github.com/jackc/pgx/v5"
)

const (
	lockRequestSQL    = "SELECT pg_advisory_xact_lock(hashtextextended($1,1))"
	findAcceptanceSQL = "SELECT order_id,request_hash FROM orders.orders WHERE request_key=$1"
	insertOrderSQL    = `INSERT INTO orders.orders(order_id,request_key,request_hash,reserve_operation_id,state,version,items) VALUES($1,$2,$3,$4,'ACCEPTED',1,$5)`
	insertJobSQL      = "INSERT INTO orders.reserve_jobs(order_id) VALUES($1)"
)

func (s *Store) Create(ctx context.Context, r model.Request) (model.Acceptance, error) {
	r, hash, err := model.Canonical(r)
	if err != nil {
		return model.Acceptance{}, err
	}
	if r.DatasetID != s.dataset {
		return model.Acceptance{}, model.ErrDataset
	}
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return model.Acceptance{}, storageError(err)
	}
	defer rollback(tx)
	if _, err = tx.Exec(ctx, lockRequestSQL, r.RequestKey); err != nil {
		return model.Acceptance{}, storageError(err)
	}
	accepted := model.Acceptance{DatasetID: s.dataset, RequestKey: r.RequestKey}
	var previousHash string
	err = tx.QueryRow(ctx, findAcceptanceSQL, r.RequestKey).Scan(&accepted.OrderID, &previousHash)
	if err == nil {
		if hash != previousHash {
			return model.Acceptance{}, model.ErrKeyConflict
		}
		return accepted, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return model.Acceptance{}, storageError(err)
	}
	var token [16]byte
	if _, err = rand.Read(token[:]); err != nil {
		return model.Acceptance{}, errors.New("cannot generate order identity")
	}
	accepted.OrderID = "order-" + hex.EncodeToString(token[:])
	items, err := json.Marshal(r.Items)
	if err != nil {
		return model.Acceptance{}, err
	}
	if _, err = tx.Exec(ctx, insertOrderSQL, accepted.OrderID, r.RequestKey, hash, "reserve-"+accepted.OrderID, items); err != nil {
		return model.Acceptance{}, storageError(err)
	}
	if _, err = tx.Exec(ctx, insertJobSQL, accepted.OrderID); err != nil {
		return model.Acceptance{}, storageError(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return model.Acceptance{}, storageError(err)
	}
	return accepted, nil
}
