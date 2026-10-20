package inventory

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	inventoryv1 "github.com/Includeoyasi/faultline-lab/api/inventory/v1"
	"github.com/Includeoyasi/faultline-lab/internal/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
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

// Runs the actual business entrypoint, TCP API and graceful signal shutdown.
// Invoked only by the isolated Linux PostgreSQL harness, including on Windows.
func TestInventoryProcess(t *testing.T) {
	s := integrationStore(t)
	fixture(t, s, "process-stock", 4, 321)
	cfg, err := database.Load(os.Getenv("FL_DATA_CONFIG"), os.Getenv("FL_DATA_SECRETS"))
	if err != nil {
		t.Fatal(err)
	}
	settings := Config{Host: cfg.Host, Port: cfg.Port, Database: "inventory_db", User: "inventory_app", SSLMode: "disable", PasswordFile: filepath.Join(cfg.Secrets, "inventory_app.password"), ConnectTimeout: "5s", RequestTimeout: "3s", MaxConnections: 2}
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "inventory.json")
	if err = os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	start := func() (inventoryv1.InventoryClient, func()) {
		t.Helper()
		cmd := exec.CommandContext(t.Context(), "/app/inventory")
		capture := &startupCapture{started: make(chan string, 1)}
		cmd.Stdout = capture
		cmd.Stderr = capture
		for _, entry := range os.Environ() {
			key, _, _ := strings.Cut(entry, "=")
			if !strings.HasPrefix(key, "FL_") {
				cmd.Env = append(cmd.Env, entry)
			}
		}
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		addressHint := listener.Addr().String()
		if err = listener.Close(); err != nil {
			t.Fatal(err)
		}
		cmd.Env = append(cmd.Env, "FL_GRPC_ADDR="+addressHint, "FL_INVENTORY_CONFIG="+path, "FL_INVENTORY_DATASET_ID="+testDataset)
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		done := make(chan error, 1)
		go func() { done <- cmd.Wait() }()
		var once sync.Once
		stop := func() {
			once.Do(func() {
				if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
					t.Error("signal Inventory process", err)
				}
				select {
				case err := <-done:
					if err != nil {
						t.Error("Inventory process exit", err)
					}
				case <-time.After(10 * time.Second):
					_ = cmd.Process.Kill()
					<-done
					t.Error("Inventory process did not drain")
				}
			})
		}
		// Always reap the owned child, including failures before startup.
		t.Cleanup(func() { once.Do(func() { _ = cmd.Process.Kill(); <-done }) })
		var address string
		select {
		case address = <-capture.started:
		case err := <-done:
			once.Do(func() {})
			t.Fatal("Inventory exited before startup", err)
		case <-time.After(10 * time.Second):
			t.Fatal("Inventory startup timed out")
		}
		conn, err := grpc.NewClient("passthrough:///"+address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = conn.Close() })
		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()
		health, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{}, grpc.WaitForReady(true))
		if err != nil || health.GetStatus() != healthpb.HealthCheckResponse_SERVING {
			t.Fatal("Inventory health", health, err)
		}
		return inventoryv1.NewInventoryClient(conn), stop
	}
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	first, stop := start()
	r := &inventoryv1.ReserveRequest{DatasetId: testDataset, OrderId: "process-order", OperationId: "process-reserve", Items: []*inventoryv1.Item{{Sku: "process-stock", Quantity: 2}}}
	original, err := first.Reserve(ctx, r)
	if err != nil || original.GetOutcome() != inventoryv1.Outcome_OUTCOME_APPLIED {
		t.Fatal(original, err)
	}
	stop()
	second, stopSecond := start()
	defer stopSecond()
	repeated, err := second.Reserve(ctx, r)
	if err != nil || !proto.Equal(original, repeated) {
		t.Fatal("process restart changed receipt", repeated, err)
	}
	committed, err := second.CommitReservation(ctx, &inventoryv1.ChangeReservationRequest{DatasetId: testDataset, OrderId: r.OrderId, OperationId: "process-commit"})
	if err != nil || committed.GetReservationAtOperation().GetState() != inventoryv1.ReservationState_RESERVATION_STATE_COMMITTED {
		t.Fatal(committed, err)
	}
	availability, err := second.GetAvailability(ctx, &inventoryv1.GetAvailabilityRequest{DatasetId: testDataset, Skus: []string{"process-stock"}})
	if err != nil || len(availability.GetItems()) != 1 || availability.Items[0].Available != 2 {
		t.Fatal("process stock", availability, err)
	}
	t.Logf("real Inventory process restart: preserved receipt, committed reservation, available=%d", availability.Items[0].Available)
}
