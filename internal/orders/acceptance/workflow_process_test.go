package acceptance

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"os"
	"testing"
	"time"

	inventoryv1 "github.com/Includeoyasi/faultline-lab/api/inventory/v1"
	ordersv1 "github.com/Includeoyasi/faultline-lab/api/orders/v1"
	"github.com/Includeoyasi/faultline-lab/internal/publishing"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type pauseAfterCommit struct {
	inventoryv1.UnimplementedInventoryServer
	peer   inventoryv1.InventoryClient
	paused chan string
}

func (p *pauseAfterCommit) Reserve(ctx context.Context, request *inventoryv1.ReserveRequest) (*inventoryv1.OperationResult, error) {
	return p.peer.Reserve(ctx, request)
}
func (p *pauseAfterCommit) GetReservation(ctx context.Context, request *inventoryv1.GetReservationRequest) (*inventoryv1.Reservation, error) {
	return p.peer.GetReservation(ctx, request)
}
func (p *pauseAfterCommit) CommitReservation(ctx context.Context, request *inventoryv1.ChangeReservationRequest) (*inventoryv1.OperationResult, error) {
	if _, err := p.peer.CommitReservation(ctx, request); err != nil {
		return nil, err
	}
	select {
	case p.paused <- request.GetOperationId():
	default:
	}
	<-ctx.Done()
	return nil, status.Error(codes.Canceled, "controlled crash after inventory commit")
}

func TestOrdersWorkflowProcessCrash(t *testing.T) {
	broker := os.Getenv("FL_PUBLISHING_TEST_BROKER")
	if broker == "" {
		t.Skip("requires isolated Kafka, PostgreSQL and Inventory")
	}
	f := setup(t)
	ctx, cancel := context.WithTimeout(t.Context(), 90*time.Second)
	defer cancel()
	client, err := kgo.NewClient(kgo.SeedBrokers(broker), kgo.ProducerBatchCompression(kgo.NoCompression()))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	request := kmsg.NewPtrCreateTopicsRequest()
	for _, name := range []string{"orders.events.v1", "payments.events.v1", "inventory.events.v1"} {
		topic := kmsg.NewCreateTopicsRequestTopic()
		topic.Topic, topic.NumPartitions, topic.ReplicationFactor = name, 3, 1
		request.Topics = append(request.Topics, topic)
	}
	response, err := request.RequestWith(ctx, client)
	if err != nil {
		t.Fatal(err)
	}
	for _, result := range response.Topics {
		if err := kerr.ErrorForCode(result.ErrorCode); err != nil && !errors.Is(err, kerr.TopicAlreadyExists) {
			t.Fatal(err)
		}
	}
	proxy := &pauseAfterCommit{peer: f.peer, paused: make(chan string, 1)}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer()
	inventoryv1.RegisterInventoryServer(server, proxy)
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)
	f.product(t, "orders-workflow-crash", 1)
	first, kill := f.startOrders(t, listener.Addr().String(), "FL_OUTBOX_BROKER="+broker)
	accepted, err := first.CreateOrder(ctx, &ordersv1.CreateOrderRequest{
		DatasetId: f.dataset, RequestKey: "orders-workflow-crash",
		Items: []*ordersv1.Item{{Sku: "orders-workflow-crash", Quantity: 1}},
	}, grpc.WaitForReady(true))
	if err != nil {
		t.Fatal(err)
	}
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	wait := func(condition func() bool) {
		t.Helper()
		for !condition() {
			select {
			case <-ctx.Done():
				t.Fatal("workflow did not reach the expected state")
			case <-ticker.C:
			}
		}
	}
	wait(func() bool { return f.read(t, accepted.OrderId).State == "PAYMENT_PENDING" })
	order := f.read(t, accepted.OrderId)
	reserved := f.inventoryEvent(t, order.OrderID, "InventoryReserved")
	if err = client.ProduceSync(ctx, reserved).FirstErr(); err != nil {
		t.Fatal(err)
	}
	wait(func() bool {
		return number(t, f.db, "SELECT count(*) FROM orders.inbox WHERE owner='inventory' AND order_id=$1 AND applied", order.OrderID) == 1
	})
	payload, err := json.Marshal(map[string]any{
		"order_id": order.OrderID, "state": "SUCCEEDED", "amount_minor": order.Reservation.TotalMinor,
		"currency": order.Reservation.Currency, "causation_id": order.OperationID + ":awaiting-payment",
	})
	if err != nil {
		t.Fatal(err)
	}
	event := publishing.Event{Dataset: f.dataset, Owner: "payments", EventID: order.OrderID + ":payment-result",
		OrderID: order.OrderID, AggregateVersion: 2, StreamSequence: 1, Type: "PaymentSucceeded",
		Payload: payload, OccurredAt: time.Now().UTC(), SchemaVersion: 1}
	body, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	if err = client.ProduceSync(ctx, &kgo.Record{Topic: "payments.events.v1", Key: []byte(order.OrderID), Value: body}).FirstErr(); err != nil {
		t.Fatal(err)
	}
	var operation string
	select {
	case operation = <-proxy.paused:
	case <-ctx.Done():
		t.Fatal("Inventory did not commit")
	}
	kill(true)
	if f.read(t, order.OrderID).State != "FINALIZING" ||
		number(t, f.inventoryDB, "SELECT count(*) FROM inventory.operations WHERE operation_id=$1", operation) != 1 {
		t.Fatal("the crash window did not retain the confirmed inventory operation")
	}
	committed := f.inventoryEvent(t, order.OrderID, "InventoryCommitted")
	if err = client.ProduceSync(ctx, committed).FirstErr(); err != nil {
		t.Fatal(err)
	}
	execute(t, f.db, "UPDATE orders.inventory_actions SET lease_until=clock_timestamp()-interval '1 second' WHERE order_id=$1", order.OrderID)
	second, stop := f.startOrders(t, listener.Addr().String(), "FL_OUTBOX_BROKER="+broker)
	defer stop(false)
	wait(func() bool {
		result, err := second.GetOrder(ctx, &ordersv1.GetOrderRequest{DatasetId: f.dataset, OrderId: order.OrderID}, grpc.WaitForReady(true))
		return err == nil && result.GetState() == ordersv1.OrderState_ORDER_STATE_COMPLETED
	})
	if number(t, f.inventoryDB, "SELECT count(*) FROM inventory.operations WHERE operation_id=$1", operation) != 1 ||
		number(t, f.db, "SELECT count(*) FROM orders.outbox WHERE order_id=$1 AND event_type='OrderCompleted'", order.OrderID) != 1 {
		t.Fatal("restart duplicated the inventory operation or terminal event")
	}
}
