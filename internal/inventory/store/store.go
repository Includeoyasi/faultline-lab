package store

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/inventory/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	bindDatasetSQL           = "INSERT INTO inventory.dataset(singleton,dataset_id) VALUES(true,$1) ON CONFLICT DO NOTHING"
	readDatasetSQL           = "SELECT dataset_id FROM inventory.dataset WHERE singleton"
	lockStockSQL             = "SELECT s.available,p.price_minor,p.currency FROM inventory.stock s JOIN inventory.products p USING(sku) WHERE s.sku=$1 FOR UPDATE OF s,p"
	insertReservationSQL     = "INSERT INTO inventory.reservations(order_id,state,version,total_minor,currency) VALUES($1,$2,1,$3,$4)"
	decreaseStockSQL         = "UPDATE inventory.stock SET available=available-$2 WHERE sku=$1"
	insertReservationItemSQL = "INSERT INTO inventory.reservation_items(order_id,sku,quantity,unit_price_minor) VALUES($1,$2,$3,$4)"
	readReservationSQL       = "SELECT state,version,total_minor,currency FROM inventory.reservations WHERE order_id=$1"
	readReservationItemsSQL  = "SELECT sku,quantity,unit_price_minor FROM inventory.reservation_items WHERE order_id=$1 ORDER BY sku"
)

// Store owns one bounded pool and exactly one immutable dataset.
type Store struct {
	pool    *pgxpool.Pool
	dataset string
	timeout time.Duration
}

// NewStore binds an already migrated database to a dataset. It never creates
// roles, schemas or tables. The caller owns the pool until NewStore succeeds.
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
	if errors.Is(err, model.ErrInvalid) {
		return model.ErrInvalid
	}
	if errors.Is(err, context.Canceled) {
		return context.Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return context.DeadlineExceeded
	}
	return errors.New("inventory storage operation failed; outcome may be unknown")
}

func (s *Store) reserve(ctx context.Context, tx pgx.Tx, r model.Request) (*model.Reservation, string, error) {
	reservation := &model.Reservation{DatasetID: s.dataset, OrderID: r.OrderID, State: model.Reserved, Version: 1, Items: []model.PricedItem{}}
	for _, item := range r.Items {
		var available, price int64
		var currency string
		err := tx.QueryRow(ctx, lockStockSQL, item.SKU).Scan(&available, &price, &currency)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.UnknownSKU, nil
		}
		if err != nil {
			return nil, "", err
		}
		if available < item.Quantity {
			return nil, model.InsufficientStock, nil
		}
		if (price > 0 && item.Quantity > math.MaxInt64/price) || reservation.TotalMinor > math.MaxInt64-price*item.Quantity {
			return nil, "", model.ErrInvalid
		}
		if reservation.Currency != "" && reservation.Currency != currency {
			return nil, "", model.ErrInvalid
		}
		reservation.Currency = currency
		reservation.TotalMinor += price * item.Quantity
		reservation.Items = append(reservation.Items, model.PricedItem{SKU: item.SKU, Quantity: item.Quantity, UnitPriceMinor: price})
	}
	if _, err := tx.Exec(ctx, insertReservationSQL, r.OrderID, model.Reserved, reservation.TotalMinor, reservation.Currency); err != nil {
		return nil, "", err
	}
	for _, item := range reservation.Items {
		if _, err := tx.Exec(ctx, decreaseStockSQL, item.SKU, item.Quantity); err != nil {
			return nil, "", err
		}
		if _, err := tx.Exec(ctx, insertReservationItemSQL, r.OrderID, item.SKU, item.Quantity, item.UnitPriceMinor); err != nil {
			return nil, "", err
		}
	}
	return reservation, model.Applied, nil
}

type querier interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

func (s *Store) readReservation(ctx context.Context, q querier, orderID string) (*model.Reservation, error) {
	r := &model.Reservation{DatasetID: s.dataset, OrderID: orderID, Items: []model.PricedItem{}}
	err := q.QueryRow(ctx, readReservationSQL, orderID).Scan(&r.State, &r.Version, &r.TotalMinor, &r.Currency)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, storageError(err)
	}
	rows, err := q.Query(ctx, readReservationItemsSQL, orderID)
	if err != nil {
		return nil, storageError(err)
	}
	defer rows.Close()
	for rows.Next() {
		var item model.PricedItem
		if err = rows.Scan(&item.SKU, &item.Quantity, &item.UnitPriceMinor); err != nil {
			return nil, storageError(err)
		}
		r.Items = append(r.Items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, storageError(err)
	}
	return r, nil
}
