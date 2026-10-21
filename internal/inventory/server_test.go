package inventory

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	inventoryv1 "github.com/Includeoyasi/faultline-lab/api/inventory/v1"
	"github.com/Includeoyasi/faultline-lab/internal/inventory/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type fakeBackend struct {
	calls int
	err   error
	block bool
}

func (b *fakeBackend) Mutate(ctx context.Context, r model.Request) (model.Receipt, error) {
	b.calls++
	if b.block {
		<-ctx.Done()
		return model.Receipt{}, ctx.Err()
	}
	return model.Receipt{OperationID: r.OperationID, Outcome: model.Applied, Reservation: &model.Reservation{DatasetID: r.DatasetID, OrderID: r.OrderID, State: model.Reserved, Version: 1}}, b.err
}
func (b *fakeBackend) GetReservation(context.Context, string, string) (*model.Reservation, error) {
	return nil, model.ErrNotFound
}
func (b *fakeBackend) GetAvailability(context.Context, string, []string) ([]model.Availability, error) {
	return nil, model.ErrNotFound
}

func TestGRPCContract(t *testing.T) {
	for _, scenario := range []string{"round trip", "invalid input", "key conflict", "dataset mismatch", "timeout", "safe storage error"} {
		t.Run(scenario, func(t *testing.T) {
			backend := &fakeBackend{}
			switch scenario {
			case "key conflict":
				backend.err = model.ErrKeyConflict
			case "dataset mismatch":
				backend.err = model.ErrDataset
			case "timeout":
				backend.block = true
			case "safe storage error":
				backend.err = errors.New("postgres://private:password@host/db")
			}
			listener := bufconn.Listen(1 << 20)
			server := grpc.NewServer()
			Register(server, backend)
			done := make(chan error, 1)
			go func() { done <- server.Serve(listener) }()
			t.Cleanup(func() {
				server.Stop()
				if err := <-done; err != nil && !errors.Is(err, grpc.ErrServerStopped) {
					t.Error(err)
				}
			})
			conn, err := grpc.NewClient("passthrough:///inventory-test", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = conn.Close() })
			ctx, cancel := context.WithTimeout(t.Context(), time.Second)
			defer cancel()
			if scenario == "timeout" {
				var shortCancel context.CancelFunc
				ctx, shortCancel = context.WithTimeout(ctx, 30*time.Millisecond)
				defer shortCancel()
			}
			r := &inventoryv1.ReserveRequest{DatasetId: "dataset", OrderId: "order", OperationId: "operation", Items: []*inventoryv1.Item{{Sku: "SKU", Quantity: 1}}}
			if scenario == "invalid input" {
				r.Items[0].Quantity = 0
			}
			response, err := inventoryv1.NewInventoryClient(conn).Reserve(ctx, r)
			want := map[string]codes.Code{"round trip": codes.OK, "invalid input": codes.InvalidArgument, "key conflict": codes.AlreadyExists, "dataset mismatch": codes.FailedPrecondition, "timeout": codes.DeadlineExceeded, "safe storage error": codes.Unavailable}[scenario]
			if status.Code(err) != want {
				t.Fatalf("got %v want %v", err, want)
			}
			if want == codes.OK && (response.GetOutcome() != inventoryv1.Outcome_OUTCOME_APPLIED || response.GetReservationAtOperation().GetState() != inventoryv1.ReservationState_RESERVATION_STATE_RESERVED) {
				t.Fatal(response)
			}
			if scenario == "invalid input" && backend.calls != 0 {
				t.Fatal("invalid request reached storage")
			}
			if scenario == "safe storage error" && status.Convert(err).Message() != "inventory storage unavailable; retry with the same operation key" {
				t.Fatal("raw storage error leaked")
			}
		})
	}
}
