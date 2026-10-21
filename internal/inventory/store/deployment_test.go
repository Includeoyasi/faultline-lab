package store

import (
	"context"
	"os"
	"testing"
	"time"

	inventoryv1 "github.com/Includeoyasi/faultline-lab/api/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

// Uses the real scratch service container, not the test binary's own Store.
// Replayed unchanged after a PostgreSQL restart to exercise pool recovery too.
func TestInventoryDeployed(t *testing.T) {
	if os.Getenv("FL_DATA_INTEGRATION") != "1" {
		t.Skip("requires isolated Inventory deployment")
	}
	target := os.Getenv("FL_INVENTORY_TARGET")
	if target == "" {
		t.Fatal("Inventory deployment target is required")
	}
	conn, err := grpc.NewClient("passthrough:///"+target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	client := inventoryv1.NewInventoryClient(conn)
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()
	r := &inventoryv1.ReserveRequest{DatasetId: testDataset, OrderId: "deployed-order", OperationId: "deployed-reserve", Items: []*inventoryv1.Item{{Sku: "SKU-001", Quantity: 1}}}
	reserved, err := client.Reserve(ctx, r)
	if err != nil || reserved.GetOutcome() != inventoryv1.Outcome_OUTCOME_APPLIED {
		t.Fatal(reserved, err)
	}
	change := &inventoryv1.ChangeReservationRequest{DatasetId: testDataset, OrderId: r.OrderId, OperationId: "deployed-release"}
	released, err := client.ReleaseReservation(ctx, change)
	if err != nil || released.GetReservationAtOperation().GetState() != inventoryv1.ReservationState_RESERVATION_STATE_RELEASED {
		t.Fatal(released, err)
	}
	repeated, err := client.ReleaseReservation(ctx, change)
	if err != nil || !proto.Equal(released, repeated) {
		t.Fatal("release RPC receipt changed", repeated, err)
	}
	current, err := client.GetReservation(ctx, &inventoryv1.GetReservationRequest{DatasetId: testDataset, OrderId: r.OrderId})
	if err != nil || current.GetState() != inventoryv1.ReservationState_RESERVATION_STATE_RELEASED {
		t.Fatal("authoritative GetReservation", current, err)
	}
	historical, err := client.Reserve(ctx, r)
	if err != nil || !proto.Equal(reserved, historical) {
		t.Fatal("old reserve RPC receipt changed", historical, err)
	}
	r.OperationId = "deployed-late-reserve"
	late, err := client.Reserve(ctx, r)
	if err != nil || late.GetOutcome() != inventoryv1.Outcome_OUTCOME_RESERVATION_EXISTS {
		t.Fatal("late reserve resurrected released stock", late, err)
	}
	t.Log("scratch Inventory container: Reserve/Release/GetReservation and immutable replay PASS")
}
