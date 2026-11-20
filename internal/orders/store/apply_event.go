package store

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	inventorymodel "github.com/Includeoyasi/faultline-lab/internal/inventory/model"
	"github.com/Includeoyasi/faultline-lab/internal/orders/model"
	"github.com/Includeoyasi/faultline-lab/internal/publishing"
	"github.com/jackc/pgx/v5"
)

func (s *Store) applyPayment(ctx context.Context, tx pgx.Tx, order model.Order, event publishing.Event) (string, error) {
	var result struct {
		OrderID     string `json:"order_id"`
		State       string `json:"state"`
		AmountMinor *int64 `json:"amount_minor"`
		Currency    string `json:"currency"`
		CausationID string `json:"causation_id"`
	}
	if json.Unmarshal(event.Payload, &result) != nil || result.OrderID != order.OrderID ||
		result.AmountMinor == nil || result.CausationID != order.OperationID+":awaiting-payment" {
		return "INVALID_PAYMENT_RESULT", nil
	}
	if order.Reservation == nil {
		return "WAITING_FOR_RESERVATION", nil
	}
	if result.Currency != order.Reservation.Currency || *result.AmountMinor != order.Reservation.TotalMinor ||
		order.State != model.PaymentPending {
		return "PAYMENT_CONTRADICTION", nil
	}
	state, action, expectedResult := model.Finalizing, "commit", "SUCCEEDED"
	if event.Type == "PaymentDeclined" {
		state, action, expectedResult = model.Cancelling, "release", "DECLINED"
	}
	if result.State != expectedResult {
		return "PAYMENT_CONTRADICTION", nil
	}
	// A stable operation key makes a timeout-after-commit safe to repeat.
	hash := sha256.Sum256([]byte(s.dataset + "\x00" + order.OrderID + "\x00" + action))
	operation := action + "-" + hex.EncodeToString(hash[:])
	if _, err := tx.Exec(ctx, `UPDATE orders.orders SET state=$2,version=version+1,updated_at=clock_timestamp()
 WHERE order_id=$1`, order.OrderID, state); err != nil {
		return "", err
	}
	_, err := tx.Exec(ctx, "INSERT INTO orders.inventory_actions(order_id,operation_id,action) VALUES($1,$2,$3)",
		order.OrderID, operation, action)
	return "", err
}

func sameReservation(order model.Order, observed *inventorymodel.Reservation) bool {
	if order.Reservation == nil || observed == nil || observed.DatasetID != order.DatasetID ||
		observed.OrderID != order.OrderID || observed.TotalMinor != order.Reservation.TotalMinor ||
		observed.Currency != order.Reservation.Currency || len(observed.Items) != len(order.Reservation.Items) {
		return false
	}
	for i, item := range observed.Items {
		expected := order.Reservation.Items[i]
		if item.SKU != expected.SKU || item.Quantity != expected.Quantity ||
			item.UnitPriceMinor != expected.UnitPriceMinor {
			return false
		}
	}
	return true
}

func (s *Store) applyInventory(ctx context.Context, tx pgx.Tx, order model.Order, event publishing.Event) (string, error) {
	var result inventorymodel.Receipt
	if json.Unmarshal(event.Payload, &result) != nil || result.Outcome != inventorymodel.Applied ||
		result.Reservation == nil || result.OperationID != event.EventID ||
		result.Reservation.Version != event.AggregateVersion {
		return "INVALID_INVENTORY_RESULT", nil
	}
	if order.Reservation == nil {
		return "WAITING_FOR_RESERVATION", nil
	}
	if !sameReservation(order, result.Reservation) {
		return "RESERVATION_CONTRADICTION", nil
	}
	if event.Type == "InventoryReserved" {
		if result.OperationID != order.OperationID || result.Reservation.State != inventorymodel.Reserved ||
			result.Reservation.Version != order.Reservation.Version {
			return "RESERVATION_CONTRADICTION", nil
		}
		return "", nil
	}
	var action, operation string
	err := tx.QueryRow(ctx, "SELECT action,operation_id FROM orders.inventory_actions WHERE order_id=$1 FOR UPDATE",
		order.OrderID).Scan(&action, &operation)
	if err == pgx.ErrNoRows {
		return "WAITING_FOR_PAYMENT_RESULT", nil
	}
	if err != nil {
		return "", err
	}
	state, expectedState, reservationState, kind := model.Completed, model.Finalizing, inventorymodel.Committed, "OrderCompleted"
	if event.Type == "InventoryReleased" {
		state, expectedState, reservationState, kind = model.Cancelled, model.Cancelling, inventorymodel.Released, "OrderCancelled"
	}
	expectedAction := "commit"
	if event.Type == "InventoryReleased" {
		expectedAction = "release"
	}
	if action != expectedAction || operation != result.OperationID || order.State != expectedState ||
		result.Reservation.State != reservationState || result.Reservation.Version != order.Reservation.Version+1 {
		return "INVENTORY_CONTRADICTION", nil
	}
	order.State = state
	order.Version++
	payload, err := json.Marshal(order)
	if err != nil {
		return "", err
	}
	if _, err = tx.Exec(ctx, `UPDATE orders.orders SET state=$2,version=$3,updated_at=clock_timestamp()
 WHERE order_id=$1`, order.OrderID, state, order.Version); err != nil {
		return "", err
	}
	if _, err = tx.Exec(ctx, "UPDATE orders.inventory_actions SET done=true,lease_until=NULL WHERE order_id=$1",
		order.OrderID); err != nil {
		return "", err
	}
	_, err = tx.Exec(ctx, `INSERT INTO orders.outbox(event_id,order_id,aggregate_version,event_type,payload)
 VALUES($1,$2,$3,$4,$5)`, order.OrderID+":"+kind, order.OrderID, order.Version, kind, payload)
	return "", err
}
