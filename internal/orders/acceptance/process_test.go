package acceptance

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	inventoryv1 "github.com/Includeoyasi/faultline-lab/api/inventory/v1"
	ordersv1 "github.com/Includeoyasi/faultline-lab/api/orders/v1"
	"github.com/Includeoyasi/faultline-lab/internal/orders"
	"github.com/Includeoyasi/faultline-lab/internal/orders/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type startupCapture struct {
	mu      sync.Mutex
	pending string
	started chan string
}

func (c *startupCapture) Write(data []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pending += string(data)
	for {
		line, rest, ok := strings.Cut(c.pending, "\n")
		if !ok {
			break
		}
		c.pending = rest
		var record struct {
			Message string `json:"msg"`
			Address string `json:"grpc_address"`
		}
		if json.Unmarshal([]byte(line), &record) == nil && record.Message == "service started" {
			select {
			case c.started <- record.Address:
			default:
			}
		}
	}
	return len(data), nil
}

type pausingInventory struct {
	inventoryv1.UnimplementedInventoryServer
	peer        inventoryv1.InventoryClient
	afterCommit bool
	first       atomic.Bool
	paused      chan *inventoryv1.ReserveRequest
}

func (p *pausingInventory) Reserve(ctx context.Context, r *inventoryv1.ReserveRequest) (*inventoryv1.OperationResult, error) {
	if !p.first.Swap(true) {
		if p.afterCommit {
			if _, err := p.peer.Reserve(ctx, r); err != nil {
				return nil, err
			}
		}
		p.paused <- proto.Clone(r).(*inventoryv1.ReserveRequest)
		<-ctx.Done()
		return nil, status.Error(codes.Canceled, "controlled lost reply")
	}
	return p.peer.Reserve(ctx, r)
}
func (p *pausingInventory) GetReservation(ctx context.Context, r *inventoryv1.GetReservationRequest) (*inventoryv1.Reservation, error) {
	return p.peer.GetReservation(ctx, r)
}
func (f *fixture) startOrders(t *testing.T, target string) (ordersv1.OrdersClient, func(bool)) {
	t.Helper()
	cfg := orders.Config{Database: store.Config{Host: f.config.Host, Port: f.config.Port, Database: "orders_db", User: "orders_app", SSLMode: f.config.SSLMode, PasswordFile: filepath.Join(f.config.Secrets, "orders_app.password"), ConnectTimeout: "5s", RequestTimeout: "3s", MaxConnections: 2}, InventoryTarget: target, WorkerCount: 1, RPCTimeout: "20s", Lease: "25s", PollInterval: "10ms", RetryMin: "10ms", RetryMax: "100ms"}
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "orders.json")
	if err = os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	hint := listener.Addr().String()
	_ = listener.Close()
	cmd := exec.CommandContext(t.Context(), "/app/orders")
	capture := &startupCapture{started: make(chan string, 1)}
	cmd.Stdout = capture
	cmd.Stderr = capture
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if !strings.HasPrefix(key, "FL_") {
			cmd.Env = append(cmd.Env, entry)
		}
	}
	cmd.Env = append(cmd.Env, "FL_GRPC_ADDR="+hint, "FL_ORDERS_CONFIG="+path, "FL_ORDERS_DATASET_ID="+f.dataset)
	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	var once sync.Once
	stop := func(hard bool) {
		once.Do(func() {
			var e error
			if hard {
				e = cmd.Process.Kill()
			} else {
				e = cmd.Process.Signal(syscall.SIGTERM)
			}
			if e != nil && !errors.Is(e, os.ErrProcessDone) {
				t.Error(e)
			}
			select {
			case e = <-done:
				if !hard && e != nil {
					t.Error("graceful Orders exit", e)
				}
			case <-time.After(10 * time.Second):
				_ = cmd.Process.Kill()
				<-done
				t.Error("Orders did not stop")
			}
		})
	}
	t.Cleanup(func() { stop(true) })
	var address string
	select {
	case address = <-capture.started:
	case err = <-done:
		once.Do(func() {})
		t.Fatal("Orders exited before startup", err)
	case <-time.After(10 * time.Second):
		t.Fatal("Orders startup timed out")
	}
	conn, err := grpc.NewClient("passthrough:///"+address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return ordersv1.NewOrdersClient(conn), stop
}
func TestOrdersProcessCrash(t *testing.T) {
	f := setup(t)
	for _, after := range []bool{false, true} {
		name := "before-reserve"
		if after {
			name = "after-inventory-commit"
		}
		t.Run(name, func(t *testing.T) {
			sku := "orders-crash-" + name
			f.product(t, sku, 2)
			proxy := &pausingInventory{peer: f.peer, afterCommit: after, paused: make(chan *inventoryv1.ReserveRequest, 1)}
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			server := grpc.NewServer()
			inventoryv1.RegisterInventoryServer(server, proxy)
			go func() { _ = server.Serve(listener) }()
			t.Cleanup(server.Stop)
			first, kill := f.startOrders(t, listener.Addr().String())
			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			request := &ordersv1.CreateOrderRequest{DatasetId: f.dataset, RequestKey: "crash-" + name, Items: []*ordersv1.Item{{Sku: sku, Quantity: 1}}}
			accepted, err := first.CreateOrder(ctx, request, grpc.WaitForReady(true))
			if err != nil {
				t.Fatal(err)
			}
			var reserve *inventoryv1.ReserveRequest
			select {
			case reserve = <-proxy.paused:
			case <-ctx.Done():
				t.Fatal("worker did not reach crash window")
			}
			kill(true)
			if number(t, f.db, "SELECT count(*) FROM orders.orders WHERE order_id=$1 AND state='RESERVING'", accepted.OrderId) != 1 {
				t.Fatal("durable attempt lost")
			}
			expected := int64(2)
			if after {
				expected = 1
			}
			if number(t, f.inventoryDB, "SELECT available FROM inventory.stock WHERE sku=$1", sku) != expected {
				t.Fatal("wrong crash window")
			}
			execute(t, f.db, "UPDATE orders.reserve_jobs SET lease_until=clock_timestamp()-interval '1 second' WHERE order_id=$1", accepted.OrderId)
			second, stop := f.startOrders(t, listener.Addr().String())
			defer stop(false)
			ticker := time.NewTicker(10 * time.Millisecond)
			defer ticker.Stop()
			for {
				result, err := second.GetOrder(ctx, &ordersv1.GetOrderRequest{DatasetId: f.dataset, OrderId: accepted.OrderId}, grpc.WaitForReady(true))
				if err != nil {
					t.Fatal(err)
				}
				if result.GetState() == ordersv1.OrderState_ORDER_STATE_PAYMENT_PENDING {
					if result.GetReserveOperationId() != reserve.GetOperationId() {
						t.Fatal("restart changed operation")
					}
					break
				}
				select {
				case <-ctx.Done():
					t.Fatal("job not resumed")
				case <-ticker.C:
				}
			}
			repeated, err := second.CreateOrder(ctx, request)
			if err != nil || !proto.Equal(accepted, repeated) {
				t.Fatal("admission changed after crash", err)
			}
			if number(t, f.inventoryDB, "SELECT available FROM inventory.stock WHERE sku=$1", sku) != 1 || number(t, f.db, "SELECT count(*) FROM orders.outbox WHERE order_id=$1", accepted.OrderId) != 1 || number(t, f.inventoryDB, "SELECT count(*) FROM inventory.operations WHERE operation_id=$1", reserve.GetOperationId()) != 1 {
				t.Fatal("crash replay duplicated effect")
			}
			t.Log("SIGKILL " + name + ": same operation, one stock effect, one local event; restarted Orders completed reserve")
		})
	}
}
