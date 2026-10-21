package store

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Includeoyasi/faultline-lab/internal/inventory/model"
	"github.com/jackc/pgx/v5"
)

const (
	lockOperationSQL          = "SELECT pg_advisory_xact_lock(1,hashtext($1))"
	readOperationSQL          = "SELECT request_hash,receipt FROM inventory.operations WHERE operation_id=$1"
	lockOrderSQL              = "SELECT pg_advisory_xact_lock(2,hashtext($1))"
	insertReleaseBarrierSQL   = "INSERT INTO inventory.reservations(order_id,state,version,total_minor,currency) VALUES($1,$2,1,0,'')"
	increaseStockSQL          = "UPDATE inventory.stock SET available=available+$2 WHERE sku=$1"
	updateReservationStateSQL = "UPDATE inventory.reservations SET state=$2,version=version+1,updated_at=now() WHERE order_id=$1"
	insertOperationSQL        = "INSERT INTO inventory.operations(operation_id,request_hash,receipt) VALUES($1,$2,$3::jsonb)"
	insertLedgerSQL           = "INSERT INTO inventory.ledger(order_id,version,sku,available_delta,operation_id) VALUES($1,$2,$3,$4,$5)"
	insertOutboxSQL           = "INSERT INTO inventory.outbox(event_id,order_id,aggregate_version,event_type,payload) VALUES($1,$2,$3,$4,$5::jsonb)"
)

// Mutate serializes operation identity AND order identity, including orders
// with no row yet. Stock rows are locked in canonical SKU order.
func (s *Store) Mutate(ctx context.Context, request model.Request) (model.Receipt, error) {
	request, hash, err := model.Canonical(request)
	if err != nil {
		return model.Receipt{}, err
	}
	if request.DatasetID != s.dataset {
		return model.Receipt{}, model.ErrDataset
	}
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return model.Receipt{}, storageError(err)
	}
	defer rollback(tx)
	// Separate advisory namespaces avoid collisions between order/key identity.
	if _, err = tx.Exec(ctx, lockOperationSQL, request.OperationID); err != nil {
		return model.Receipt{}, storageError(err)
	}
	var previousHash string
	var data []byte
	err = tx.QueryRow(ctx, readOperationSQL, request.OperationID).Scan(&previousHash, &data)
	if err == nil {
		if previousHash != hash {
			return model.Receipt{}, model.ErrKeyConflict
		}
		var result model.Receipt
		if err = json.Unmarshal(data, &result); err != nil {
			return model.Receipt{}, storageError(err)
		}
		return result, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return model.Receipt{}, storageError(err)
	}
	if _, err = tx.Exec(ctx, lockOrderSQL, request.OrderID); err != nil {
		return model.Receipt{}, storageError(err)
	}
	reservation, err := s.readReservation(ctx, tx, request.OrderID)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return model.Receipt{}, err
	}
	err = nil // An absent order is an expected input to the state machine.
	result := model.Receipt{OperationID: request.OperationID, Outcome: model.Applied, Reservation: reservation}
	changed := false
	switch request.Method {
	case "reserve":
		if reservation != nil {
			result.Outcome = model.ReservationExists
			break
		}
		result.Reservation, result.Outcome, err = s.reserve(ctx, tx, request)
		changed = err == nil && result.Outcome == model.Applied
	case "commit", "release":
		target := model.Committed
		if request.Method == "release" {
			target = model.Released
		}
		switch {
		case reservation == nil && target == model.Released:
			reservation = &model.Reservation{DatasetID: s.dataset, OrderID: request.OrderID, State: model.Released, Version: 1, Items: []model.PricedItem{}}
			_, err = tx.Exec(ctx, insertReleaseBarrierSQL, request.OrderID, model.Released)
			result.Reservation = reservation
			changed = err == nil
		case reservation == nil:
			result.Outcome = model.InvalidState
		case reservation.State == target: // Already at requested terminal state: no second effect.
		case reservation.State != model.Reserved:
			result.Outcome = model.InvalidState
		default:
			if target == model.Released {
				for _, item := range reservation.Items {
					if _, err = tx.Exec(ctx, increaseStockSQL, item.SKU, item.Quantity); err != nil {
						break
					}
				}
			}
			if err == nil {
				_, err = tx.Exec(ctx, updateReservationStateSQL, request.OrderID, target)
			}
			if err == nil {
				reservation.State = target
				reservation.Version++
				changed = true
			}
		}
	}
	if err != nil {
		return model.Receipt{}, storageError(err)
	}
	data, err = json.Marshal(result)
	if err != nil {
		return model.Receipt{}, storageError(err)
	}
	if _, err = tx.Exec(ctx, insertOperationSQL, request.OperationID, hash, string(data)); err != nil {
		return model.Receipt{}, storageError(err)
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
			if _, err = tx.Exec(ctx, insertLedgerSQL, request.OrderID, result.Reservation.Version, item.SKU, delta, request.OperationID); err != nil {
				return model.Receipt{}, storageError(err)
			}
		}
		eventType := map[string]string{"reserve": "InventoryReserved", "commit": "InventoryCommitted", "release": "InventoryReleased"}[request.Method]
		if _, err = tx.Exec(ctx, insertOutboxSQL, request.OperationID, request.OrderID, result.Reservation.Version, eventType, string(data)); err != nil {
			return model.Receipt{}, storageError(err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return model.Receipt{}, storageError(err)
	}
	return result, nil
}
