package payments

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	inventoryv1 "github.com/Includeoyasi/faultline-lab/api/inventory/v1"
	ordersv1 "github.com/Includeoyasi/faultline-lab/api/orders/v1"
	"github.com/Includeoyasi/faultline-lab/internal/inventory"
	inventorystore "github.com/Includeoyasi/faultline-lab/internal/inventory/store"
	"github.com/Includeoyasi/faultline-lab/internal/orders"
	orderstore "github.com/Includeoyasi/faultline-lab/internal/orders/store"
	"github.com/Includeoyasi/faultline-lab/internal/orders/worker"
	"github.com/Includeoyasi/faultline-lab/internal/orders/workflow"
	"github.com/Includeoyasi/faultline-lab/internal/publishing"
	"github.com/Includeoyasi/faultline-lab/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type lostCommitReply struct {
	inventoryv1.InventoryClient
	calls chan string
	first atomic.Bool
}

func (c *lostCommitReply) CommitReservation(ctx context.Context, request *inventoryv1.ChangeReservationRequest, options ...grpc.CallOption) (*inventoryv1.OperationResult, error) {
	result, err := c.InventoryClient.CommitReservation(ctx, request, options...)
	if err == nil {
		select {
		case c.calls <- request.GetOperationId():
		default:
		}
		if !c.first.Swap(true) {
			return nil, status.Error(codes.DeadlineExceeded, "controlled lost response after commit")
		}
	}
	return result, err
}

// TestOrderPaymentWorkflow uses the real service handlers, SQL stores,
// Kafka publishers and consumers. The database and broker are disposable.
func TestOrderPaymentWorkflow(t *testing.T) {
	if os.Getenv("FL_RUN_ORDER_PAYMENT_WORKFLOW") != "1" {
		t.Skip("requires isolated workflow broker and PostgreSQL")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 140*time.Second)
	defer cancel()
	check := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	dsn, broker := os.Getenv("FL_PUBLISHING_TEST_DB"), os.Getenv("FL_PUBLISHING_TEST_BROKER")
	admin, err := pgxpool.New(ctx, dsn)
	check(err)
	defer admin.Close()
	_, err = admin.Exec(ctx, "CREATE DATABASE order_payment_workflow_test")
	check(err)
	cfg, err := pgxpool.ParseConfig(dsn)
	check(err)
	cfg.ConnConfig.Database = "order_payment_workflow_test"
	cfg.MaxConns = 20
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	check(err)
	defer pool.Close()
	for _, owner := range []string{"orders", "inventory", "payments"} {
		paths, err := fs.Glob(migrations.Files, owner+"/*.sql")
		check(err)
		for _, path := range paths {
			data, err := migrations.Files.ReadFile(path)
			check(err)
			_, err = pool.Exec(ctx, string(data))
			check(err)
		}
	}
	const dataset = "workflow-e2e"
	stock, err := inventorystore.NewStore(ctx, pool, dataset, 5*time.Second)
	check(err)
	backend, err := orderstore.NewStore(ctx, pool, dataset, 5*time.Second)
	check(err)
	payment := processor{pool: pool, dataset: dataset, declineBPS: 5000, maxDelay: time.Millisecond}
	check(payment.bind(ctx))
	_, err = pool.Exec(ctx, `INSERT INTO inventory.products(sku,title,price_minor,currency) VALUES('workflow-sku','Workflow',1250,'RUB');
 INSERT INTO inventory.stock(sku,available) VALUES('workflow-sku',2)`)
	check(err)
	listen := func(register func(*grpc.Server)) string {
		t.Helper()
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		check(err)
		server := grpc.NewServer()
		register(server)
		go func() { _ = server.Serve(listener) }()
		t.Cleanup(server.Stop)
		return listener.Addr().String()
	}
	inventoryAddress := listen(func(s *grpc.Server) { inventory.Register(s, stock) })
	orderAddress := listen(func(s *grpc.Server) { orders.Register(s, backend) })
	dial := func(address string) *grpc.ClientConn {
		t.Helper()
		conn, err := grpc.NewClient("passthrough:///"+address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		t.Cleanup(func() { _ = conn.Close() })
		return conn
	}
	inventoryClient := inventoryv1.NewInventoryClient(dial(inventoryAddress))
	orderClient := ordersv1.NewOrdersClient(dial(orderAddress))
	var ids [2]string
	for i := range ids {
		request := &ordersv1.CreateOrderRequest{DatasetId: dataset, RequestKey: fmt.Sprintf("workflow-%d", i),
			Items: []*ordersv1.Item{{Sku: "workflow-sku", Quantity: 1}}}
		result, err := orderClient.CreateOrder(ctx, request, grpc.WaitForReady(true))
		check(err)
		ids[i] = result.GetOrderId()
	}
	// Keep exactly two orders while exercising both deterministic simulator outcomes.
	found := false
	for n := 0; n < 1000; n++ {
		payment.seed = fmt.Sprintf("workflow-seed-%d", n)
		first, _ := payment.simulate(ids[0])
		second, _ := payment.simulate(ids[1])
		if first != second {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("no simulator seed produced both outcomes")
	}
	commitProxy := &lostCommitReply{InventoryClient: inventoryClient, calls: make(chan string, 16)}
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	reserve := worker.Processor{Jobs: backend, Inventory: inventoryClient, Lease: 8 * time.Second,
		RPCTimeout: 2 * time.Second, RetryMin: 100 * time.Millisecond, RetryMax: time.Second, Logger: logger}
	advance := workflow.Processor{Backend: backend, Inventory: commitProxy, Lease: 8 * time.Second,
		RPCTimeout: 2 * time.Second, RetryMin: 100 * time.Millisecond, RetryMax: time.Second, Logger: logger}
	runCtx, stop := context.WithCancel(ctx)
	var wg sync.WaitGroup
	start := func(fn func()) { wg.Add(1); go func() { defer wg.Done(); fn() }() }
	defer func() { stop(); wg.Wait() }()
	start(func() { reserve.Run(runCtx, 2, 20*time.Millisecond) })
	start(func() {
		publishing.Publisher{Pool: pool, Owner: "orders", Dataset: dataset, Broker: broker}.Run(runCtx, logger)
	})
	start(func() {
		publishing.Publisher{Pool: pool, Owner: "payments", Dataset: dataset, Broker: broker}.Run(runCtx, logger)
	})
	consumer, err := newConsumer(broker, dataset, "orders.events.v1")
	check(err)
	defer func() { stop(); consumer.Close() }()
	consumerError := make(chan error, 1)
	start(func() { consumerError <- payment.receive(runCtx, consumer) })
	start(func() { payment.work(runCtx, logger) })
	startWorkflow := func() (context.CancelFunc, <-chan error) {
		workerCtx, workerStop := context.WithCancel(runCtx)
		done := make(chan error, 1)
		start(func() { done <- advance.Run(workerCtx, broker, dataset, 2, 20*time.Millisecond) })
		return workerStop, done
	}
	stopWorkflow, workflowDone := startWorkflow()
	var attempts [2]string
	for i := range attempts {
		select {
		case attempts[i] = <-commitProxy.calls:
		case err := <-consumerError:
			t.Fatalf("payment consumer stopped: %v", err)
		case err := <-workflowDone:
			t.Fatalf("order workflow stopped: %v", err)
		case <-ctx.Done():
			t.Fatal("commit was not retried")
		}
	}
	if attempts[0] == "" || attempts[0] != attempts[1] {
		t.Fatalf("operation identity changed after lost response: %v", attempts)
	}
	// Stop and restart Orders after the Inventory mutation but before its
	// confirmation event can be published.
	stopWorkflow()
	if err := <-workflowDone; err != nil && !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	_, workflowDone = startWorkflow()
	start(func() {
		publishing.Publisher{Pool: pool, Owner: "inventory", Dataset: dataset, Broker: broker}.Run(runCtx, logger)
	})

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	completed, cancelled := "", ""
	for completed == "" || cancelled == "" {
		select {
		case <-ctx.Done():
			t.Fatal("orders did not reach terminal states")
		case err := <-workflowDone:
			t.Fatalf("order workflow stopped: %v", err)
		case err := <-consumerError:
			t.Fatalf("payment consumer stopped: %v", err)
		case <-ticker.C:
			for _, id := range ids {
				current, err := orderClient.GetOrder(ctx, &ordersv1.GetOrderRequest{DatasetId: dataset, OrderId: id})
				check(err)
				switch current.GetState() {
				case ordersv1.OrderState_ORDER_STATE_COMPLETED:
					completed = id
				case ordersv1.OrderState_ORDER_STATE_CANCELLED:
					cancelled = id
				}
			}
		}
	}
	count := func(query string, args ...any) int64 {
		t.Helper()
		var n int64
		check(pool.QueryRow(ctx, query, args...).Scan(&n))
		return n
	}
	if count("SELECT available FROM inventory.stock WHERE sku='workflow-sku'") != 1 ||
		count("SELECT count(*) FROM payments.debits WHERE order_id=$1 AND amount_minor=1250 AND currency='RUB'", completed) != 1 ||
		count("SELECT count(*) FROM payments.debits WHERE order_id=$1", cancelled) != 0 ||
		count("SELECT count(*) FROM orders.outbox WHERE event_type IN ('OrderCompleted','OrderCancelled')") != 2 {
		t.Fatal("stock, payment or order outcomes disagree")
	}
	if count("SELECT count(*) FROM inventory.operations WHERE operation_id=$1", attempts[0]) != 1 {
		t.Fatal("the repeated commit had more than one effect")
	}
	// Identical replay is harmless, but a conflicting identity or missing
	// sequence cannot advance the durable inbox.
	var raw []byte
	check(pool.QueryRow(ctx, "SELECT body FROM orders.inbox WHERE owner='payments' AND order_id=$1", completed).Scan(&raw))
	record := &kgo.Record{Topic: "payments.events.v1", Key: []byte(completed), Value: raw}
	check(backend.ReceiveEvent(ctx, record))
	var event publishing.Event
	check(json.Unmarshal(raw, &event))
	event.Type = "PaymentDeclined"
	record.Value, err = json.Marshal(event)
	check(err)
	if backend.ReceiveEvent(ctx, record) == nil {
		t.Fatal("conflicting event accepted")
	}
	event.EventID += "-gap"
	event.StreamSequence = 3
	record.Value, err = json.Marshal(event)
	check(err)
	if backend.ReceiveEvent(ctx, record) == nil {
		t.Fatal("missing sequence accepted")
	}
}
