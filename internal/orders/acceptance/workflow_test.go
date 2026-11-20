package acceptance

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/orders/model"
	"github.com/Includeoyasi/faultline-lab/internal/orders/workflow"
	"github.com/Includeoyasi/faultline-lab/internal/publishing"
	"github.com/twmb/franz-go/pkg/kgo"
)

// This complements the Kafka workflow test by using the actual limited
// orders_app and inventory_app identities from the PostgreSQL smoke.
func testPaymentWorkflow(t *testing.T, f *fixture) {
	for _, scenario := range []struct {
		name, event, paymentState, middle, terminal, inventoryEvent, orderEvent string
		stock                                                                   int64
	}{
		{"success", "PaymentSucceeded", "SUCCEEDED", model.Finalizing, model.Completed, "InventoryCommitted", "OrderCompleted", 0},
		{"declined", "PaymentDeclined", "DECLINED", model.Cancelling, model.Cancelled, "InventoryReleased", "OrderCancelled", 1},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			sku := "orders-workflow-" + scenario.name
			f.product(t, sku, 1)
			accepted := f.create(t, sku, sku)
			f.step(t)
			order := f.read(t, accepted.OrderID)
			if order.State != model.PaymentPending || order.Reservation == nil {
				t.Fatalf("reservation not saved: %+v", order)
			}
			inventoryRecord := func(kind string) *kgo.Record {
				t.Helper()
				var event publishing.Event
				event.Dataset, event.Owner, event.OrderID, event.SchemaVersion = f.dataset, "inventory", order.OrderID, 1
				err := f.inventoryDB.QueryRow(t.Context(), `SELECT event_id,aggregate_version,stream_sequence,event_type,payload,created_at
 FROM inventory.outbox WHERE order_id=$1 AND event_type=$2`, order.OrderID, kind).
					Scan(&event.EventID, &event.AggregateVersion, &event.StreamSequence, &event.Type, &event.Payload, &event.OccurredAt)
				if err != nil {
					t.Fatal(err)
				}
				data, err := json.Marshal(event)
				if err != nil {
					t.Fatal(err)
				}
				return &kgo.Record{Topic: "inventory.events.v1", Key: []byte(order.OrderID), Value: data}
			}
			reserved := inventoryRecord("InventoryReserved")
			if err := f.store.ReceiveEvent(t.Context(), reserved); err != nil {
				t.Fatal(err)
			}
			if worked, err := f.store.ApplyEvent(t.Context()); err != nil || !worked {
				t.Fatal(worked, err)
			}
			payload, err := json.Marshal(map[string]any{
				"order_id": order.OrderID, "state": scenario.paymentState,
				"amount_minor": order.Reservation.TotalMinor, "currency": order.Reservation.Currency,
				"causation_id": order.OperationID + ":awaiting-payment",
			})
			if err != nil {
				t.Fatal(err)
			}
			event := publishing.Event{Dataset: f.dataset, Owner: "payments", OrderID: order.OrderID,
				EventID: order.OrderID + ":payment-result", AggregateVersion: 2, StreamSequence: 1,
				Type: scenario.event, Payload: payload, SchemaVersion: 1, OccurredAt: time.Now().UTC()}
			data, err := json.Marshal(event)
			if err != nil {
				t.Fatal(err)
			}
			payment := &kgo.Record{Topic: "payments.events.v1", Key: []byte(order.OrderID), Value: data}
			if err = f.store.ReceiveEvent(t.Context(), payment); err != nil {
				t.Fatal(err)
			}
			if worked, err := f.store.ApplyEvent(t.Context()); err != nil || !worked {
				t.Fatal(worked, err)
			}
			if f.read(t, order.OrderID).State != scenario.middle {
				t.Fatal("payment did not create the expected workflow state")
			}
			p := workflow.Processor{Backend: f.store, Inventory: f.peer, Lease: 30 * time.Second,
				RPCTimeout: 5 * time.Second, RetryMin: time.Millisecond, RetryMax: time.Second}
			if worked, err := p.Step(t.Context()); err != nil || !worked {
				t.Fatal(worked, err)
			}
			if f.read(t, order.OrderID).State != scenario.middle {
				t.Fatal("RPC reply alone prematurely completed the order")
			}
			terminal := inventoryRecord(scenario.inventoryEvent)
			if err = f.store.ReceiveEvent(t.Context(), terminal); err != nil {
				t.Fatal(err)
			}
			if worked, err := f.store.ApplyEvent(t.Context()); err != nil || !worked {
				t.Fatal(worked, err)
			}
			if f.read(t, order.OrderID).State != scenario.terminal ||
				number(t, f.inventoryDB, "SELECT available FROM inventory.stock WHERE sku=$1", sku) != scenario.stock {
				t.Fatal("terminal state or stock disagrees with the outcome")
			}
			for _, record := range []*kgo.Record{payment, reserved, terminal} {
				if err = f.store.ReceiveEvent(t.Context(), record); err != nil {
					t.Fatal(err)
				}
			}
			if number(t, f.db, "SELECT count(*) FROM orders.outbox WHERE order_id=$1 AND event_type=$2", order.OrderID, scenario.orderEvent) != 1 {
				t.Fatal("terminal result was not recorded exactly once")
			}
		})
	}
}
