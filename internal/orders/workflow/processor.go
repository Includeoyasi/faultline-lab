// Package workflow advances Orders after a payment result.
package workflow

import (
	"context"
	"errors"
	"log/slog"
	"math/rand/v2"
	"sync"
	"time"

	inventoryv1 "github.com/Includeoyasi/faultline-lab/api/inventory/v1"
	"github.com/Includeoyasi/faultline-lab/internal/orders/model"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Backend interface {
	ReceiveEvent(context.Context, *kgo.Record) error
	ApplyEvent(context.Context) (bool, error)
	ClaimInventoryAction(context.Context, time.Duration) (*model.InventoryAction, error)
	RecordInventoryAttempt(context.Context, model.InventoryAction, string, string, time.Duration) error
}

type Inventory interface {
	CommitReservation(context.Context, *inventoryv1.ChangeReservationRequest, ...grpc.CallOption) (*inventoryv1.OperationResult, error)
	ReleaseReservation(context.Context, *inventoryv1.ChangeReservationRequest, ...grpc.CallOption) (*inventoryv1.OperationResult, error)
}

type Processor struct {
	Backend    Backend
	Inventory  Inventory
	Lease      time.Duration
	RPCTimeout time.Duration
	RetryMin   time.Duration
	RetryMax   time.Duration
	Logger     *slog.Logger
}

// Run joins every local worker when the receiver exits or the service stops.
func (p *Processor) Run(ctx context.Context, broker, dataset string, count int, poll time.Duration) error {
	client, err := kgo.NewClient(kgo.SeedBrokers(broker),
		kgo.ConsumeTopics("payments.events.v1", "inventory.events.v1"),
		kgo.ConsumerGroup("orders-workflow-v1-"+dataset),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
		kgo.FetchIsolationLevel(kgo.ReadCommitted()), kgo.DisableAutoCommit(), kgo.BlockRebalanceOnPoll())
	if err != nil {
		return err
	}
	defer client.Close()
	workers, stop := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Go(func() { p.applyLoop(workers, poll) })
	for range count {
		wg.Go(func() { p.actionLoop(workers, poll) })
	}
	err = p.receive(workers, client)
	stop()
	wg.Wait()
	return err
}

func (p *Processor) receive(ctx context.Context, client *kgo.Client) error {
	defer client.AllowRebalance()
	for ctx.Err() == nil {
		fetch := client.PollRecords(ctx, 1)
		if errs := fetch.Errors(); len(errs) > 0 {
			return errs[0].Err
		}
		for _, record := range fetch.Records() {
			if err := p.Backend.ReceiveEvent(ctx, record); err != nil {
				return err
			}
			if err := client.CommitRecords(ctx, record); err != nil {
				return err
			}
		}
		client.AllowRebalance()
	}
	return ctx.Err()
}

func (p *Processor) applyLoop(ctx context.Context, poll time.Duration) {
	for ctx.Err() == nil {
		worked, err := p.Backend.ApplyEvent(ctx)
		if err != nil && ctx.Err() == nil && p.Logger != nil {
			p.Logger.WarnContext(ctx, "order event application delayed", "error", err)
		}
		if (!worked || err != nil) && !sleep(ctx, poll) {
			return
		}
	}
}

func (p *Processor) actionLoop(ctx context.Context, poll time.Duration) {
	for ctx.Err() == nil {
		worked, err := p.Step(ctx)
		if err != nil && ctx.Err() == nil && p.Logger != nil {
			p.Logger.WarnContext(ctx, "inventory action delayed", "error", err)
		}
		if (!worked || err != nil) && !sleep(ctx, poll) {
			return
		}
	}
}

func sleep(ctx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

// Step sends one idempotent operation and retains the job until the
// corresponding Inventory event is applied.
func (p *Processor) Step(ctx context.Context) (bool, error) {
	job, err := p.Backend.ClaimInventoryAction(ctx, p.Lease)
	if err != nil || job == nil {
		return false, err
	}
	rpcCtx, cancel := context.WithTimeout(ctx, p.RPCTimeout)
	request := &inventoryv1.ChangeReservationRequest{
		DatasetId: job.Order.DatasetID, OrderId: job.Order.OrderID, OperationId: job.OperationID,
	}
	var result *inventoryv1.OperationResult
	if job.Action == "commit" {
		result, err = p.Inventory.CommitReservation(rpcCtx, request)
	} else {
		result, err = p.Inventory.ReleaseReservation(rpcCtx, request)
	}
	cancel()
	blocked, reason := "", ""
	if err != nil {
		switch status.Code(err) {
		case codes.DeadlineExceeded, codes.Canceled, codes.Unavailable, codes.ResourceExhausted, codes.Aborted, codes.Internal, codes.Unknown:
			reason = "INVENTORY_UNCONFIRMED"
		default:
			blocked = "INVENTORY_CONTRACT_ERROR"
		}
	} else if !validResult(*job, result) {
		blocked = "INVENTORY_CONTRADICTION"
	}
	delay := p.RetryMin
	for n := int64(1); n < job.Attempts && delay < p.RetryMax; n++ {
		if delay > p.RetryMax/2 {
			delay = p.RetryMax
			break
		}
		delay *= 2
	}
	if delay > p.RetryMin {
		delay = p.RetryMin + time.Duration(rand.Int64N(int64(delay-p.RetryMin)+1))
	}
	err = p.Backend.RecordInventoryAttempt(ctx, *job, blocked, reason, delay)
	if errors.Is(err, model.ErrStale) {
		return true, nil
	}
	return true, err
}

func validResult(job model.InventoryAction, result *inventoryv1.OperationResult) bool {
	expected := job.Order.Reservation
	if expected == nil || result == nil || result.GetOperationId() != job.OperationID ||
		result.GetOutcome() != inventoryv1.Outcome_OUTCOME_APPLIED {
		return false
	}
	reservation := result.GetReservationAtOperation()
	state := inventoryv1.ReservationState_RESERVATION_STATE_COMMITTED
	if job.Action == "release" {
		state = inventoryv1.ReservationState_RESERVATION_STATE_RELEASED
	}
	if reservation == nil || reservation.GetDatasetId() != job.Order.DatasetID ||
		reservation.GetOrderId() != job.Order.OrderID || reservation.GetState() != state ||
		reservation.GetVersion() != expected.Version+1 || reservation.GetTotalMinor() != expected.TotalMinor ||
		reservation.GetCurrency() != expected.Currency || len(reservation.GetItems()) != len(expected.Items) {
		return false
	}
	for i, item := range reservation.GetItems() {
		want := expected.Items[i]
		if item == nil || item.GetSku() != want.SKU || item.GetQuantity() != want.Quantity ||
			item.GetUnitPriceMinor() != want.UnitPriceMinor {
			return false
		}
	}
	return true
}
