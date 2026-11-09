// Package worker executes the one durable reserve step owned by Orders.
package worker

import (
	"context"
	"errors"
	"log/slog"
	"math"
	"math/rand/v2"
	"reflect"
	"sync"
	"time"

	inventoryv1 "github.com/Includeoyasi/faultline-lab/api/inventory/v1"
	"github.com/Includeoyasi/faultline-lab/internal/orders/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Jobs interface {
	Claim(context.Context, time.Duration) (*model.Job, error)
	Resolve(context.Context, model.Job, model.Resolution, time.Duration) error
}
type Inventory interface {
	Reserve(context.Context, *inventoryv1.ReserveRequest, ...grpc.CallOption) (*inventoryv1.OperationResult, error)
	GetReservation(context.Context, *inventoryv1.GetReservationRequest, ...grpc.CallOption) (*inventoryv1.Reservation, error)
}
type Processor struct {
	Jobs       Jobs
	Inventory  Inventory
	Lease      time.Duration
	RPCTimeout time.Duration
	RetryMin   time.Duration
	RetryMax   time.Duration
	Logger     *slog.Logger
}

func (p *Processor) Step(ctx context.Context) (bool, error) {
	job, err := p.Jobs.Claim(ctx, p.Lease)
	if err != nil || job == nil {
		return false, err
	}
	rpcCtx, cancel := context.WithTimeout(ctx, p.RPCTimeout)
	resolution := p.reserve(rpcCtx, job.Order)
	cancel()
	delay := p.RetryMin
	for attempt := int64(1); attempt < job.Attempts && delay < p.RetryMax; attempt++ {
		if delay > p.RetryMax/2 {
			delay = p.RetryMax
			break
		}
		delay *= 2
	}
	if delay > p.RetryMin {
		delay = p.RetryMin + time.Duration(rand.Int64N(int64(delay-p.RetryMin)+1))
	}
	err = p.Jobs.Resolve(ctx, *job, resolution, delay)
	if errors.Is(err, model.ErrStale) {
		return true, nil
	}
	if err == nil && p.Logger != nil {
		p.Logger.InfoContext(ctx, "reserve attempt recorded", "order_id", job.Order.OrderID, "operation_id", job.Order.OperationID, "generation", job.Generation, "attempt", job.Attempts, "blocked_reason", resolution.BlockedReason, "retry_reason", resolution.RetryReason, "rejection", resolution.Rejection)
	}
	return true, err
}

// Run has bounded parallelism and joins every worker before returning.
func (p *Processor) Run(ctx context.Context, count int, poll time.Duration) {
	var wg sync.WaitGroup
	for range count {
		wg.Go(func() {
			timer := time.NewTimer(0)
			defer timer.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-timer.C:
				}
				worked, err := p.Step(ctx)
				if err != nil && ctx.Err() == nil && p.Logger != nil {
					p.Logger.ErrorContext(ctx, "reserve attempt unavailable; durable job retained")
				}
				wait := poll
				if worked && err == nil {
					wait = 0
				}
				timer.Reset(wait)
			}
		})
	}
	wg.Wait()
}
func (p *Processor) reserve(ctx context.Context, o model.Order) model.Resolution {
	request := &inventoryv1.ReserveRequest{DatasetId: o.DatasetID, OrderId: o.OrderID, OperationId: o.OperationID}
	for _, item := range o.Items {
		request.Items = append(request.Items, &inventoryv1.Item{Sku: item.SKU, Quantity: item.Quantity})
	}
	result, err := p.Inventory.Reserve(ctx, request)
	if err != nil {
		return peerError(err)
	}
	if result == nil || result.GetOperationId() != o.OperationID {
		return model.Resolution{BlockedReason: "INVALID_RESERVE_IDENTITY"}
	}
	switch result.GetOutcome() {
	case inventoryv1.Outcome_OUTCOME_INSUFFICIENT_STOCK, inventoryv1.Outcome_OUTCOME_UNKNOWN_SKU:
		if result.GetReservationAtOperation() != nil {
			return model.Resolution{BlockedReason: "INVALID_REJECTION_SNAPSHOT"}
		}
		reason := "UNKNOWN_SKU"
		if result.GetOutcome() == inventoryv1.Outcome_OUTCOME_INSUFFICIENT_STOCK {
			reason = "INSUFFICIENT_STOCK"
		}
		return model.Resolution{Rejection: reason}
	case inventoryv1.Outcome_OUTCOME_APPLIED:
		historical, ok := verified(o, result.GetReservationAtOperation())
		if !ok {
			return model.Resolution{BlockedReason: "INVALID_RESERVE_SNAPSHOT"}
		}
		current, err := p.Inventory.GetReservation(ctx, &inventoryv1.GetReservationRequest{DatasetId: o.DatasetID, OrderId: o.OrderID})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return model.Resolution{RetryReason: "RESERVATION_NOT_YET_OBSERVED"}
			}
			return peerError(err)
		}
		observed, ok := verified(o, current)
		if !ok || !reflect.DeepEqual(historical, observed) {
			return model.Resolution{BlockedReason: "RESERVATION_CONTRADICTION"}
		}
		return model.Resolution{Reservation: observed}
	default:
		return model.Resolution{BlockedReason: "UNEXPECTED_RESERVE_OUTCOME"}
	}
}
func peerError(err error) model.Resolution {
	switch status.Code(err) {
	case codes.DeadlineExceeded, codes.Canceled, codes.Unavailable, codes.ResourceExhausted, codes.Aborted, codes.Internal, codes.Unknown:
		return model.Resolution{RetryReason: "INVENTORY_UNCONFIRMED"}
	default:
		return model.Resolution{BlockedReason: "INVENTORY_CONTRACT_ERROR"}
	}
}
func verified(o model.Order, r *inventoryv1.Reservation) (*model.Reservation, bool) {
	if r == nil || r.GetDatasetId() != o.DatasetID || r.GetOrderId() != o.OrderID || r.GetState() != inventoryv1.ReservationState_RESERVATION_STATE_RESERVED || r.GetVersion() < 1 || len(r.GetItems()) != len(o.Items) || len(r.GetCurrency()) != 3 {
		return nil, false
	}
	for _, c := range r.GetCurrency() {
		if c < 'A' || c > 'Z' {
			return nil, false
		}
	}
	bySKU := make(map[string]*inventoryv1.PricedItem, len(r.GetItems()))
	for _, item := range r.GetItems() {
		if item == nil || bySKU[item.GetSku()] != nil {
			return nil, false
		}
		bySKU[item.GetSku()] = item
	}
	snapshot := &model.Reservation{Version: r.GetVersion(), Currency: r.GetCurrency(), Items: []model.PricedItem{}}
	for _, item := range o.Items {
		priced := bySKU[item.SKU]
		if priced == nil || priced.GetQuantity() != item.Quantity || priced.GetUnitPriceMinor() < 0 {
			return nil, false
		}
		price := priced.GetUnitPriceMinor()
		if price > 0 && item.Quantity > math.MaxInt64/price {
			return nil, false
		}
		amount := item.Quantity * price
		if snapshot.TotalMinor > math.MaxInt64-amount {
			return nil, false
		}
		snapshot.TotalMinor += amount
		snapshot.Items = append(snapshot.Items, model.PricedItem{SKU: item.SKU, Quantity: item.Quantity, UnitPriceMinor: price})
	}
	return snapshot, snapshot.TotalMinor == r.GetTotalMinor()
}
