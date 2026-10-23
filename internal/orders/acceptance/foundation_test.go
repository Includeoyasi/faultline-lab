// Package acceptance exercises Orders with real PostgreSQL and Inventory gRPC.
package acceptance

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	inventoryv1 "github.com/Includeoyasi/faultline-lab/api/inventory/v1"
	ordersv1 "github.com/Includeoyasi/faultline-lab/api/orders/v1"
	"github.com/Includeoyasi/faultline-lab/internal/database"
	"github.com/Includeoyasi/faultline-lab/internal/orders"
	"github.com/Includeoyasi/faultline-lab/internal/orders/model"
	"github.com/Includeoyasi/faultline-lab/internal/orders/store"
	"github.com/Includeoyasi/faultline-lab/internal/orders/worker"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type fixture struct {
	store       *store.Store
	db          *pgx.Conn
	inventoryDB *pgx.Conn
	config      database.Config
	dataset     string
	peer        inventoryv1.InventoryClient
}

func setup(t *testing.T) *fixture {
	t.Helper()
	if os.Getenv("FL_DATA_INTEGRATION") != "1" {
		t.Skip("requires isolated PostgreSQL and Inventory")
	}
	cfg, err := database.Load(os.Getenv("FL_DATA_CONFIG"), os.Getenv("FL_DATA_SECRETS"))
	if err != nil {
		t.Fatal(err)
	}
	dataset := os.Getenv("FL_ORDERS_DATASET_ID")
	if dataset == "" {
		t.Fatal("Orders dataset required")
	}
	backend, err := store.Open(t.Context(), store.Config{Host: cfg.Host, Port: cfg.Port, Database: "orders_db", User: "orders_app", SSLMode: cfg.SSLMode, PasswordFile: filepath.Join(cfg.Secrets, "orders_app.password"), DatasetID: dataset, ConnectTimeout: "5s", RequestTimeout: "10s", MaxConnections: 4})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(backend.Close)
	db, err := database.Connect(t.Context(), cfg, "orders_db", "orders_app")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close(db) })
	inv, err := database.Connect(t.Context(), cfg, "inventory_db", "inventory_app")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close(inv) })
	target := os.Getenv("FL_INVENTORY_TARGET")
	if target == "" {
		t.Fatal("Inventory target required")
	}
	conn, err := grpc.NewClient("passthrough:///"+target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return &fixture{store: backend, db: db, inventoryDB: inv, config: cfg, dataset: dataset, peer: inventoryv1.NewInventoryClient(conn)}
}
func execute(t *testing.T, db *pgx.Conn, sql string, args ...any) {
	t.Helper()
	if _, err := db.Exec(t.Context(), sql, args...); err != nil {
		t.Fatal(err)
	}
}
func number(t *testing.T, db *pgx.Conn, sql string, args ...any) int64 {
	t.Helper()
	var n int64
	if err := db.QueryRow(t.Context(), sql, args...).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}
func (f *fixture) product(t *testing.T, sku string, quantity int64) {
	t.Helper()
	execute(t, f.inventoryDB, "INSERT INTO inventory.products(sku,title,price_minor,currency) VALUES($1,$1,123,'RUB')", sku)
	execute(t, f.inventoryDB, "INSERT INTO inventory.stock(sku,available) VALUES($1,$2)", sku, quantity)
}
func (f *fixture) request(key, sku string) model.Request {
	return model.Request{DatasetID: f.dataset, RequestKey: key, Items: []model.Item{{SKU: sku, Quantity: 1}}}
}
func (f *fixture) processor(peer worker.Inventory) *worker.Processor {
	return &worker.Processor{Jobs: f.store, Inventory: peer, Lease: 30 * time.Second, RPCTimeout: 5 * time.Second, RetryMin: time.Millisecond, RetryMax: time.Second}
}
func (f *fixture) create(t *testing.T, key, sku string) model.Acceptance {
	t.Helper()
	r, err := f.store.Create(t.Context(), f.request(key, sku))
	if err != nil {
		t.Fatal(err)
	}
	return r
}
func (f *fixture) read(t *testing.T, id string) model.Order {
	t.Helper()
	o, err := f.store.GetOrder(t.Context(), f.dataset, id)
	if err != nil {
		t.Fatal(err)
	}
	return o
}
func (f *fixture) step(t *testing.T) {
	t.Helper()
	worked, err := f.processor(f.peer).Step(t.Context())
	if err != nil || !worked {
		t.Fatal(worked, err)
	}
}
func (f *fixture) ddlFault(t *testing.T, table string) func() {
	t.Helper()
	migrator, err := database.Connect(t.Context(), f.config, "orders_db", "orders_migrator")
	if err != nil {
		t.Fatal(err)
	}
	// Only fixed, test-owned table identifiers are permitted here.
	if table != "reserve_jobs" && table != "outbox" {
		t.Fatal("invalid fault target")
	}
	execute(t, migrator, "ALTER TABLE orders."+table+" ADD CONSTRAINT test_insert_failure CHECK(false) NOT VALID")
	var once sync.Once
	restore := func() {
		once.Do(func() {
			execute(t, migrator, "ALTER TABLE orders."+table+" DROP CONSTRAINT test_insert_failure")
			database.Close(migrator)
		})
	}
	t.Cleanup(restore)
	return restore
}
func TestOrdersFoundation(t *testing.T) {
	f := setup(t)
	t.Run("concurrent admission immutable receipt and current read", func(t *testing.T) {
		f.product(t, "orders-admission", 5)
		request := f.request("admission", "orders-admission")
		accepted := make(chan model.Acceptance, 20)
		var wg sync.WaitGroup
		for range 20 {
			wg.Go(func() {
				r, err := f.store.Create(t.Context(), request)
				if err != nil {
					t.Error(err)
					return
				}
				accepted <- r
			})
		}
		wg.Wait()
		close(accepted)
		var first model.Acceptance
		for r := range accepted {
			if first.OrderID == "" {
				first = r
			}
			if r != first {
				t.Fatal("different admission receipt")
			}
		}
		if first.OrderID == "" || number(t, f.db, "SELECT count(*) FROM orders.reserve_jobs WHERE order_id=$1", first.OrderID) != 1 {
			t.Fatal("missing unique job")
		}
		changed := request
		changed.Items = []model.Item{{SKU: "orders-admission", Quantity: 2}}
		if _, err := f.store.Create(t.Context(), changed); !errors.Is(err, model.ErrKeyConflict) {
			t.Fatal(err)
		}
		f.step(t)
		if f.read(t, first.OrderID).State != model.PaymentPending {
			t.Fatal("reserve did not progress")
		}
		repeated, err := f.store.Create(t.Context(), request)
		if err != nil || repeated != first {
			t.Fatal("progress changed acceptance", err)
		}
		current, err := f.store.GetByRequestKey(t.Context(), f.dataset, "admission")
		if err != nil || current.State != model.PaymentPending || current.Reservation.TotalMinor != 123 {
			t.Fatal(current, err)
		}
		if number(t, f.inventoryDB, "SELECT available FROM inventory.stock WHERE sku='orders-admission'") != 4 || number(t, f.db, "SELECT count(*) FROM orders.outbox WHERE order_id=$1", first.OrderID) != 1 {
			t.Fatal("duplicate effect")
		}
		if _, err = f.store.GetOrder(t.Context(), "other", first.OrderID); !errors.Is(err, model.ErrDataset) {
			t.Fatal("foreign dataset accepted", err)
		}
	})
	t.Run("failed job insert rolls back admission", func(t *testing.T) {
		restore := f.ddlFault(t, "reserve_jobs")
		if _, err := f.store.Create(t.Context(), f.request("failed-admission", "SKU-001")); err == nil {
			t.Fatal("injected failure ignored")
		}
		if number(t, f.db, "SELECT count(*) FROM orders.orders WHERE request_key='failed-admission'") != 0 {
			t.Fatal("partial admission")
		}
		restore()
	})
	t.Run("outbox failure rolls back result and replay recovers", func(t *testing.T) {
		f.product(t, "orders-outbox", 2)
		a := f.create(t, "outbox", "orders-outbox")
		restore := f.ddlFault(t, "outbox")
		if _, err := f.processor(f.peer).Step(t.Context()); err == nil {
			t.Fatal("injected event failure ignored")
		}
		o := f.read(t, a.OrderID)
		if o.State != model.Reserving || o.Reservation != nil || number(t, f.db, "SELECT count(*) FROM orders.reserve_jobs WHERE order_id=$1 AND done", a.OrderID) != 0 {
			t.Fatal("partial completion")
		}
		restore()
		execute(t, f.db, "UPDATE orders.reserve_jobs SET lease_until=clock_timestamp()-interval '1 second' WHERE order_id=$1", a.OrderID)
		f.step(t)
		if f.read(t, a.OrderID).State != model.PaymentPending || number(t, f.inventoryDB, "SELECT available FROM inventory.stock WHERE sku='orders-outbox'") != 1 || number(t, f.db, "SELECT count(*) FROM orders.outbox WHERE order_id=$1", a.OrderID) != 1 {
			t.Fatal("replay duplicated effect")
		}
	})
	t.Run("expired old worker cannot finish or reschedule", func(t *testing.T) {
		a := f.create(t, "takeover", "absent-sku")
		old, err := f.store.Claim(t.Context(), 30*time.Second)
		if err != nil || old == nil {
			t.Fatal(err)
		}
		execute(t, f.db, "UPDATE orders.reserve_jobs SET lease_until=clock_timestamp()-interval '1 second' WHERE order_id=$1", a.OrderID)
		fresh, err := f.store.Claim(t.Context(), 30*time.Second)
		if err != nil || fresh == nil || fresh.Generation <= old.Generation || fresh.Order.OperationID != old.Order.OperationID {
			t.Fatal(fresh, err)
		}
		for _, r := range []model.Resolution{{RetryReason: "OLD"}, {BlockedReason: "OLD"}, {Reservation: &model.Reservation{Version: 1, Items: []model.PricedItem{{SKU: "absent-sku", Quantity: 1, UnitPriceMinor: 1}}, TotalMinor: 1, Currency: "RUB"}}} {
			if err = f.store.Resolve(t.Context(), *old, r, 0); !errors.Is(err, model.ErrStale) {
				t.Fatal("stale effect accepted", err)
			}
		}
		if err = f.store.Resolve(t.Context(), *fresh, model.Resolution{Rejection: "UNKNOWN_SKU"}, 0); err != nil {
			t.Fatal(err)
		}
		if f.read(t, a.OrderID).State != model.Rejected || number(t, f.db, "SELECT count(*) FROM orders.outbox WHERE order_id=$1", a.OrderID) != 0 {
			t.Fatal("old result survived takeover")
		}
	})
	t.Run("lease expiry during order lock wait rolls back the result", func(t *testing.T) {
		a := f.create(t, "lease-wait", "lease-wait-sku")
		job, err := f.store.Claim(t.Context(), 30*time.Second)
		if err != nil || job == nil {
			t.Fatal(err)
		}
		lock, err := f.db.Begin(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = lock.Rollback(context.Background()) }()
		if _, err = lock.Exec(t.Context(), "SELECT 1 FROM orders.orders WHERE order_id=$1 FOR UPDATE", a.OrderID); err != nil {
			t.Fatal(err)
		}
		observer, err := database.Connect(t.Context(), f.config, "orders_db", "orders_app")
		if err != nil {
			t.Fatal(err)
		}
		defer database.Close(observer)
		execute(t, observer, "UPDATE orders.reserve_jobs SET lease_until=clock_timestamp()+interval '300 milliseconds' WHERE order_id=$1", a.OrderID)
		done := make(chan error, 1)
		go func() {
			done <- f.store.Resolve(t.Context(), *job, model.Resolution{Reservation: &model.Reservation{Version: 1, Items: []model.PricedItem{{SKU: "lease-wait-sku", Quantity: 1, UnitPriceMinor: 1}}, TotalMinor: 1, Currency: "RUB"}}, 0)
		}()
		ticker := time.NewTicker(5 * time.Millisecond)
		defer ticker.Stop()
		deadline := time.NewTimer(3 * time.Second)
		defer deadline.Stop()
		for number(t, observer, "SELECT count(*) FROM pg_stat_activity WHERE usename='orders_app' AND wait_event_type='Lock' AND query LIKE 'SELECT o.order_id%'") == 0 {
			select {
			case <-deadline.C:
				t.Fatal("Resolve did not reach the locked order")
			case err := <-done:
				t.Fatal("Resolve exited before the lock wait", err)
			case <-ticker.C:
			}
		}
		for number(t, observer, "SELECT count(*) FROM orders.reserve_jobs WHERE order_id=$1 AND lease_until<=clock_timestamp()", a.OrderID) == 0 {
			select {
			case <-deadline.C:
				t.Fatal("database lease did not expire")
			case <-ticker.C:
			}
		}
		if err = lock.Rollback(t.Context()); err != nil {
			t.Fatal(err)
		}
		select {
		case err = <-done:
			if !errors.Is(err, model.ErrStale) {
				t.Fatal("expired result accepted", err)
			}
		case <-time.After(3 * time.Second):
			t.Fatal("result did not finish")
		}
		if f.read(t, a.OrderID).State != model.Reserving || number(t, f.db, "SELECT count(*) FROM orders.outbox WHERE order_id=$1", a.OrderID) != 0 {
			t.Fatal("expired result leaked writes")
		}
		fresh, err := f.store.Claim(t.Context(), 30*time.Second)
		if err != nil || fresh == nil {
			t.Fatal(err)
		}
		if err = f.store.Resolve(t.Context(), *fresh, model.Resolution{Rejection: "UNKNOWN_SKU"}, 0); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("historical success with released current reservation stays unresolved", func(t *testing.T) {
		f.product(t, "orders-contradiction", 2)
		a := f.create(t, "contradiction", "orders-contradiction")
		o := f.read(t, a.OrderID)
		request := &inventoryv1.ReserveRequest{DatasetId: f.dataset, OrderId: o.OrderID, OperationId: o.OperationID, Items: []*inventoryv1.Item{{Sku: "orders-contradiction", Quantity: 1}}}
		if _, err := f.peer.Reserve(t.Context(), request); err != nil {
			t.Fatal(err)
		}
		if _, err := f.peer.ReleaseReservation(t.Context(), &inventoryv1.ChangeReservationRequest{DatasetId: f.dataset, OrderId: o.OrderID, OperationId: "release-contradiction"}); err != nil {
			t.Fatal(err)
		}
		f.step(t)
		current := f.read(t, a.OrderID)
		if current.State != model.Reserving || current.BlockedReason != "RESERVATION_CONTRADICTION" || current.Reservation != nil || number(t, f.db, "SELECT count(*) FROM orders.outbox WHERE order_id=$1", a.OrderID) != 0 {
			t.Fatal("false success from historical receipt", current)
		}
	})
	t.Run("1000 accepted attempts compete for 100 units", func(t *testing.T) {
		f.product(t, "orders-g1", 100)
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		server := grpc.NewServer()
		orders.Register(server, f.store)
		go func() { _ = server.Serve(listener) }()
		t.Cleanup(server.Stop)
		conn, err := grpc.NewClient("passthrough:///"+listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = conn.Close() }()
		client := ordersv1.NewOrdersClient(conn)
		workerCtx, stop := context.WithCancel(t.Context())
		done := make(chan struct{})
		go func() { f.processor(f.peer).Run(workerCtx, 4, time.Millisecond); close(done) }()
		defer func() { stop(); <-done }()
		ids := make(chan string, 1000)
		slots := make(chan struct{}, 8)
		var wg sync.WaitGroup
		for i := range 1000 {
			slots <- struct{}{}
			wg.Go(func() {
				defer func() { <-slots }()
				ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
				defer cancel()
				r, err := client.CreateOrder(ctx, &ordersv1.CreateOrderRequest{DatasetId: f.dataset, RequestKey: fmt.Sprintf("g1-%d", i), Items: []*ordersv1.Item{{Sku: "orders-g1", Quantity: 1}}})
				if err != nil {
					t.Error(err)
					return
				}
				ids <- r.GetOrderId()
			})
		}
		wg.Wait()
		close(ids)
		deadline := time.NewTimer(60 * time.Second)
		defer deadline.Stop()
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		for number(t, f.db, "SELECT count(*) FROM orders.orders WHERE request_key LIKE 'g1-%' AND state IN ('PAYMENT_PENDING','REJECTED')") != 1000 {
			select {
			case <-deadline.C:
				t.Fatal("accepted jobs did not settle")
			case <-ticker.C:
			}
		}
		if len(ids) != 1000 || number(t, f.db, "SELECT count(*) FROM orders.orders WHERE request_key LIKE 'g1-%' AND state='PAYMENT_PENDING'") != 100 || number(t, f.db, "SELECT count(*) FROM orders.orders WHERE request_key LIKE 'g1-%' AND rejection='INSUFFICIENT_STOCK'") != 900 || number(t, f.inventoryDB, "SELECT available FROM inventory.stock WHERE sku='orders-g1'") != 0 || number(t, f.inventoryDB, "SELECT sum(quantity) FROM inventory.reservation_items WHERE sku='orders-g1'") != 100 {
			t.Fatal("G1 conservation failed")
		}
		t.Log("1000 attempts: 100 PAYMENT_PENDING, 900 durable stock rejections, stock balance 0")
	})
}
func TestOrdersDeployed(t *testing.T) {
	if os.Getenv("FL_DATA_INTEGRATION") != "1" {
		t.Skip("requires business containers")
	}
	conn, err := grpc.NewClient("passthrough:///"+os.Getenv("FL_ORDERS_TARGET"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	client := ordersv1.NewOrdersClient(conn)
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()
	dataset := os.Getenv("FL_ORDERS_DATASET_ID")
	request := &ordersv1.CreateOrderRequest{DatasetId: dataset, RequestKey: "deployed-orders", Items: []*ordersv1.Item{{Sku: "SKU-001", Quantity: 1}}}
	first, err := client.CreateOrder(ctx, request, grpc.WaitForReady(true))
	if err != nil || first.GetState() != ordersv1.OrderState_ORDER_STATE_ACCEPTED {
		t.Fatal(first, err)
	}
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	for {
		current, err := client.GetOrderByRequestKey(ctx, &ordersv1.GetOrderByRequestKeyRequest{DatasetId: dataset, RequestKey: request.RequestKey})
		if err != nil {
			t.Fatal(err)
		}
		if current.GetState() == ordersv1.OrderState_ORDER_STATE_PAYMENT_PENDING {
			if current.GetTotalMinor() <= 0 || current.GetBlockedReason() != "" {
				t.Fatal(current)
			}
			break
		}
		select {
		case <-ctx.Done():
			t.Fatal("deployment did not reserve")
		case <-ticker.C:
		}
	}
	repeated, err := client.CreateOrder(ctx, request)
	if err != nil || !proto.Equal(first, repeated) {
		t.Fatal("deployment receipt changed", err)
	}
	read, err := client.GetOrder(ctx, &ordersv1.GetOrderRequest{DatasetId: dataset, OrderId: first.OrderId})
	if err != nil || read.GetState() != ordersv1.OrderState_ORDER_STATE_PAYMENT_PENDING {
		t.Fatal(read, err)
	}
	request.Items[0].Quantity++
	if _, err = client.CreateOrder(ctx, request); status.Code(err) != codes.AlreadyExists {
		t.Fatal("changed payload accepted", err)
	}
	t.Log("scratch Orders: immutable admission, durable reserve and current reads PASS")
}
