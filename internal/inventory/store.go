package inventory

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	if !identifier.MatchString(dataset) || timeout < time.Millisecond || timeout > time.Minute {
		return nil, ErrInvalid
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, storageError(err)
	}
	defer rollback(tx)
	if _, err = tx.Exec(ctx, "INSERT INTO inventory.dataset(singleton,dataset_id) VALUES(true,$1) ON CONFLICT DO NOTHING", dataset); err != nil {
		return nil, storageError(err)
	}
	var actual string
	if err = tx.QueryRow(ctx, "SELECT dataset_id FROM inventory.dataset WHERE singleton").Scan(&actual); err != nil {
		return nil, storageError(err)
	}
	if actual != dataset {
		return nil, ErrDataset
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, storageError(err)
	}
	return &Store{pool: pool, dataset: dataset, timeout: timeout}, nil
}

func (s *Store) Close() { s.pool.Close() }

func rollback(tx pgx.Tx) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = tx.Rollback(ctx)
}
func storageError(err error) error {
	if errors.Is(err, ErrInvalid) {
		return ErrInvalid
	}
	if errors.Is(err, context.Canceled) {
		return context.Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return context.DeadlineExceeded
	}
	return errors.New("inventory storage operation failed; outcome may be unknown")
}

// Mutate serializes operation identity AND order identity, including orders
// with no row yet. Stock rows are locked in canonical SKU order.
func (s *Store) Mutate(ctx context.Context, request Request) (Receipt, error) {
	request, hash, err := canonical(request)
	if err != nil {
		return Receipt{}, err
	}
	if request.DatasetID != s.dataset {
		return Receipt{}, ErrDataset
	}
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Receipt{}, storageError(err)
	}
	defer rollback(tx)
	// Separate advisory namespaces avoid collisions between order/key identity.
	if _, err = tx.Exec(ctx, "SELECT pg_advisory_xact_lock(1,hashtext($1))", request.OperationID); err != nil {
		return Receipt{}, storageError(err)
	}
	var previousHash string
	var data []byte
	err = tx.QueryRow(ctx, "SELECT request_hash,receipt FROM inventory.operations WHERE operation_id=$1", request.OperationID).Scan(&previousHash, &data)
	if err == nil {
		if previousHash != hash {
			return Receipt{}, ErrKeyConflict
		}
		var result Receipt
		if err = json.Unmarshal(data, &result); err != nil {
			return Receipt{}, storageError(err)
		}
		return result, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return Receipt{}, storageError(err)
	}
	if _, err = tx.Exec(ctx, "SELECT pg_advisory_xact_lock(2,hashtext($1))", request.OrderID); err != nil {
		return Receipt{}, storageError(err)
	}
	reservation, err := s.readReservation(ctx, tx, request.OrderID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return Receipt{}, err
	}
	err = nil // An absent order is an expected input to the state machine.
	result := Receipt{OperationID: request.OperationID, Outcome: Applied, Reservation: reservation}
	changed := false
	switch request.Method {
	case "reserve":
		if reservation != nil {
			result.Outcome = ReservationExists
			break
		}
		result.Reservation, result.Outcome, err = s.reserve(ctx, tx, request)
		changed = err == nil && result.Outcome == Applied
	case "commit", "release":
		target := Committed
		if request.Method == "release" {
			target = Released
		}
		switch {
		case reservation == nil && target == Released:
			reservation = &Reservation{DatasetID: s.dataset, OrderID: request.OrderID, State: Released, Version: 1, Items: []PricedItem{}}
			_, err = tx.Exec(ctx, "INSERT INTO inventory.reservations(order_id,state,version,total_minor,currency) VALUES($1,$2,1,0,'')", request.OrderID, Released)
			result.Reservation = reservation
			changed = err == nil
		case reservation == nil:
			result.Outcome = InvalidState
		case reservation.State == target: // Already at requested terminal state: no second effect.
		case reservation.State != Reserved:
			result.Outcome = InvalidState
		default:
			if target == Released {
				for _, item := range reservation.Items {
					if _, err = tx.Exec(ctx, "UPDATE inventory.stock SET available=available+$2 WHERE sku=$1", item.SKU, item.Quantity); err != nil {
						break
					}
				}
			}
			if err == nil {
				_, err = tx.Exec(ctx, "UPDATE inventory.reservations SET state=$2,version=version+1,updated_at=now() WHERE order_id=$1", request.OrderID, target)
			}
			if err == nil {
				reservation.State = target
				reservation.Version++
				changed = true
			}
		}
	}
	if err != nil {
		return Receipt{}, storageError(err)
	}
	data, err = json.Marshal(result)
	if err != nil {
		return Receipt{}, storageError(err)
	}
	if _, err = tx.Exec(ctx, "INSERT INTO inventory.operations(operation_id,request_hash,receipt) VALUES($1,$2,$3::jsonb)", request.OperationID, hash, string(data)); err != nil {
		return Receipt{}, storageError(err)
	}
	if changed {
		for _, item := range result.Reservation.Items {
			delta := int64(0)
			if request.Method == "reserve" {
				delta = -item.Quantity
			}
			if request.Method == "release" {
				delta = item.Quantity
			}
			if _, err = tx.Exec(ctx, "INSERT INTO inventory.ledger(order_id,version,sku,available_delta,operation_id) VALUES($1,$2,$3,$4,$5)", request.OrderID, result.Reservation.Version, item.SKU, delta, request.OperationID); err != nil {
				return Receipt{}, storageError(err)
			}
		}
		eventType := map[string]string{"reserve": "InventoryReserved", "commit": "InventoryCommitted", "release": "InventoryReleased"}[request.Method]
		if _, err = tx.Exec(ctx, "INSERT INTO inventory.outbox(event_id,order_id,aggregate_version,event_type,payload) VALUES($1,$2,$3,$4,$5::jsonb)", request.OperationID, request.OrderID, result.Reservation.Version, eventType, string(data)); err != nil {
			return Receipt{}, storageError(err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return Receipt{}, storageError(err)
	}
	return result, nil
}

func (s *Store) reserve(ctx context.Context, tx pgx.Tx, r Request) (*Reservation, string, error) {
	reservation := &Reservation{DatasetID: s.dataset, OrderID: r.OrderID, State: Reserved, Version: 1, Items: []PricedItem{}}
	for _, item := range r.Items {
		var available, price int64
		var currency string
		err := tx.QueryRow(ctx, "SELECT s.available,p.price_minor,p.currency FROM inventory.stock s JOIN inventory.products p USING(sku) WHERE s.sku=$1 FOR UPDATE OF s,p", item.SKU).Scan(&available, &price, &currency)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, UnknownSKU, nil
		}
		if err != nil {
			return nil, "", err
		}
		if available < item.Quantity {
			return nil, InsufficientStock, nil
		}
		if (price > 0 && item.Quantity > math.MaxInt64/price) || reservation.TotalMinor > math.MaxInt64-price*item.Quantity {
			return nil, "", ErrInvalid
		}
		if reservation.Currency != "" && reservation.Currency != currency {
			return nil, "", ErrInvalid
		}
		reservation.Currency = currency
		reservation.TotalMinor += price * item.Quantity
		reservation.Items = append(reservation.Items, PricedItem{SKU: item.SKU, Quantity: item.Quantity, UnitPriceMinor: price})
	}
	if _, err := tx.Exec(ctx, "INSERT INTO inventory.reservations(order_id,state,version,total_minor,currency) VALUES($1,$2,1,$3,$4)", r.OrderID, Reserved, reservation.TotalMinor, reservation.Currency); err != nil {
		return nil, "", err
	}
	for _, item := range reservation.Items {
		if _, err := tx.Exec(ctx, "UPDATE inventory.stock SET available=available-$2 WHERE sku=$1", item.SKU, item.Quantity); err != nil {
			return nil, "", err
		}
		if _, err := tx.Exec(ctx, "INSERT INTO inventory.reservation_items(order_id,sku,quantity,unit_price_minor) VALUES($1,$2,$3,$4)", r.OrderID, item.SKU, item.Quantity, item.UnitPriceMinor); err != nil {
			return nil, "", err
		}
	}
	return reservation, Applied, nil
}

type querier interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

func (s *Store) readReservation(ctx context.Context, q querier, orderID string) (*Reservation, error) {
	r := &Reservation{DatasetID: s.dataset, OrderID: orderID, Items: []PricedItem{}}
	err := q.QueryRow(ctx, "SELECT state,version,total_minor,currency FROM inventory.reservations WHERE order_id=$1", orderID).Scan(&r.State, &r.Version, &r.TotalMinor, &r.Currency)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, storageError(err)
	}
	rows, err := q.Query(ctx, "SELECT sku,quantity,unit_price_minor FROM inventory.reservation_items WHERE order_id=$1 ORDER BY sku", orderID)
	if err != nil {
		return nil, storageError(err)
	}
	defer rows.Close()
	for rows.Next() {
		var item PricedItem
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

func (s *Store) GetReservation(ctx context.Context, dataset, orderID string) (*Reservation, error) {
	if dataset != s.dataset {
		return nil, ErrDataset
	}
	if !identifier.MatchString(orderID) {
		return nil, ErrInvalid
	}
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.readReservation(ctx, s.pool, orderID)
}

func (s *Store) GetAvailability(ctx context.Context, dataset string, skus []string) ([]Availability, error) {
	if dataset != s.dataset {
		return nil, ErrDataset
	}
	if len(skus) == 0 || len(skus) > 100 {
		return nil, ErrInvalid
	}
	seen := map[string]bool{}
	for _, sku := range skus {
		if !skuPattern.MatchString(sku) || seen[sku] {
			return nil, ErrInvalid
		}
		seen[sku] = true
	}
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	rows, err := s.pool.Query(ctx, "SELECT s.sku,s.available,p.price_minor,p.currency FROM inventory.stock s JOIN inventory.products p USING(sku) WHERE s.sku=ANY($1) ORDER BY s.sku", skus)
	if err != nil {
		return nil, storageError(err)
	}
	defer rows.Close()
	result := []Availability{}
	for rows.Next() {
		var item Availability
		if err = rows.Scan(&item.SKU, &item.Available, &item.UnitPriceMinor, &item.Currency); err != nil {
			return nil, storageError(err)
		}
		result = append(result, item)
	}
	if err = rows.Err(); err != nil {
		return nil, storageError(err)
	}
	if len(result) != len(skus) {
		return nil, ErrNotFound
	}
	return result, nil
}
