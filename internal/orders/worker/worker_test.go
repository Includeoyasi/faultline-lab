package worker

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	inventoryv1 "github.com/Includeoyasi/faultline-lab/api/inventory/v1"
	"github.com/Includeoyasi/faultline-lab/internal/orders/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type fakeJobs struct {
	job         model.Job
	resolution  model.Resolution
	delay       time.Duration
	finishError error
}

func (f *fakeJobs) Claim(context.Context, time.Duration) (*model.Job, error) { return &f.job, nil }
func (f *fakeJobs) Resolve(_ context.Context, _ model.Job, r model.Resolution, d time.Duration) error {
	f.resolution = r
	f.delay = d
	return f.finishError
}

type fakeInventory struct {
	result    *inventoryv1.OperationResult
	current   *inventoryv1.Reservation
	err       error
	readErr   error
	operation string
}

func (f *fakeInventory) Reserve(_ context.Context, r *inventoryv1.ReserveRequest, _ ...grpc.CallOption) (*inventoryv1.OperationResult, error) {
	f.operation = r.OperationId
	return f.result, f.err
}
func (f *fakeInventory) GetReservation(context.Context, *inventoryv1.GetReservationRequest, ...grpc.CallOption) (*inventoryv1.Reservation, error) {
	return f.current, f.readErr
}
func TestReserveOutcomeGuards(t *testing.T) {
	order := model.Order{DatasetID: "data", OrderID: "order", OperationID: "reserve", State: model.Reserving, Items: []model.Item{{SKU: "sku", Quantity: 2}}}
	good := &inventoryv1.Reservation{DatasetId: "data", OrderId: "order", State: inventoryv1.ReservationState_RESERVATION_STATE_RESERVED, Version: 1, Items: []*inventoryv1.PricedItem{{Sku: "sku", Quantity: 2, UnitPriceMinor: 100}}, TotalMinor: 200, Currency: "RUB"}
	for _, name := range []string{"success", "timeout", "refused", "unknown-sku", "unexpected", "wrong-key", "wrong-dataset", "wrong-order", "wrong-items", "wrong-total", "negative-price", "overflow", "duplicate-sku", "wrong-currency", "missing", "released", "committed", "version-changed", "not-found", "lookup-unavailable", "stale-finish"} {
		t.Run(name, func(t *testing.T) {
			historical := proto.Clone(good).(*inventoryv1.Reservation)
			current := proto.Clone(good).(*inventoryv1.Reservation)
			peer := &fakeInventory{result: &inventoryv1.OperationResult{OperationId: "reserve", Outcome: inventoryv1.Outcome_OUTCOME_APPLIED, ReservationAtOperation: historical}, current: current}
			jobs := &fakeJobs{job: model.Job{Order: order, Generation: 1, OrderVersion: 2, Attempts: 50}}
			want := "blocked"
			switch name {
			case "success":
				want = "success"
			case "timeout":
				peer.err = status.Error(codes.DeadlineExceeded, "unknown")
				want = "retry"
			case "refused":
				peer.result.Outcome = inventoryv1.Outcome_OUTCOME_INSUFFICIENT_STOCK
				peer.result.ReservationAtOperation = nil
				want = "reject"
			case "unknown-sku":
				peer.result.Outcome = inventoryv1.Outcome_OUTCOME_UNKNOWN_SKU
				peer.result.ReservationAtOperation = nil
				want = "reject"
			case "unexpected":
				peer.result.Outcome = inventoryv1.Outcome_OUTCOME_RESERVATION_EXISTS
			case "wrong-key":
				peer.result.OperationId = "other"
			case "wrong-dataset":
				historical.DatasetId = "other"
			case "wrong-order":
				historical.OrderId = "other"
			case "wrong-items":
				historical.Items[0].Quantity = 3
			case "wrong-total":
				historical.TotalMinor++
			case "negative-price":
				historical.Items[0].UnitPriceMinor = -1
			case "overflow":
				historical.Items[0].UnitPriceMinor = math.MaxInt64
			case "duplicate-sku":
				historical.Items = append(historical.Items, proto.Clone(historical.Items[0]).(*inventoryv1.PricedItem))
			case "wrong-currency":
				historical.Currency = "rub"
			case "missing":
				peer.result.ReservationAtOperation = nil
			case "released":
				current.State = inventoryv1.ReservationState_RESERVATION_STATE_RELEASED
			case "committed":
				current.State = inventoryv1.ReservationState_RESERVATION_STATE_COMMITTED
			case "version-changed":
				current.Version++
			case "not-found":
				peer.readErr = status.Error(codes.NotFound, "not yet visible")
				want = "retry"
			case "lookup-unavailable":
				peer.readErr = status.Error(codes.Unavailable, "offline")
				want = "retry"
			case "stale-finish":
				jobs.finishError = model.ErrStale
				want = "success"
			}
			p := Processor{Jobs: jobs, Inventory: peer, Lease: time.Second, RPCTimeout: time.Second, RetryMin: time.Millisecond, RetryMax: time.Second}
			if worked, err := p.Step(t.Context()); !worked || err != nil {
				t.Fatal(worked, err)
			}
			if peer.operation != order.OperationID {
				t.Fatal("changed operation identity")
			}
			r := jobs.resolution
			switch want {
			case "success":
				if r.Reservation == nil {
					t.Fatal(r)
				}
			case "retry":
				if r.RetryReason == "" || jobs.delay < time.Millisecond || jobs.delay > time.Second {
					t.Fatal(r, jobs.delay)
				}
			case "reject":
				if r.Rejection == "" {
					t.Fatal(r)
				}
			case "blocked":
				if r.BlockedReason == "" || r.Reservation != nil {
					t.Fatal(r)
				}
			}
		})
	}
}
func TestCanceledWorkerKeepsAttemptUnconfirmed(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	jobs := &fakeJobs{job: model.Job{Order: model.Order{OperationID: "stable"}}, finishError: context.Canceled}
	p := Processor{Jobs: jobs, Inventory: &fakeInventory{err: status.Error(codes.Canceled, "canceled")}, RPCTimeout: time.Second, RetryMin: time.Millisecond, RetryMax: time.Second}
	if _, err := p.Step(ctx); !errors.Is(err, context.Canceled) || jobs.resolution.Rejection != "" {
		t.Fatal(err, jobs.resolution)
	}
}
