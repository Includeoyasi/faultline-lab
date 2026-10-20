package inventory

import (
	"context"
	"errors"

	inventoryv1 "github.com/Includeoyasi/faultline-lab/api/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Backend defines the domain boundary; SQL correctness is tested on PostgreSQL.
type Backend interface {
	Mutate(context.Context, Request) (Receipt, error)
	GetReservation(context.Context, string, string) (*Reservation, error)
	GetAvailability(context.Context, string, []string) ([]Availability, error)
}
type Server struct {
	inventoryv1.UnimplementedInventoryServer
	backend Backend
}

func Register(server *grpc.Server, backend Backend) {
	inventoryv1.RegisterInventoryServer(server, &Server{backend: backend})
}

func (s *Server) Reserve(ctx context.Context, r *inventoryv1.ReserveRequest) (*inventoryv1.OperationResult, error) {
	request := Request{DatasetID: r.GetDatasetId(), OrderID: r.GetOrderId(), OperationID: r.GetOperationId(), Method: "reserve"}
	for _, item := range r.GetItems() {
		request.Items = append(request.Items, Item{SKU: item.GetSku(), Quantity: item.GetQuantity()})
	}
	return s.mutate(ctx, request)
}
func (s *Server) CommitReservation(ctx context.Context, r *inventoryv1.ChangeReservationRequest) (*inventoryv1.OperationResult, error) {
	return s.mutate(ctx, Request{DatasetID: r.GetDatasetId(), OrderID: r.GetOrderId(), OperationID: r.GetOperationId(), Method: "commit"})
}
func (s *Server) ReleaseReservation(ctx context.Context, r *inventoryv1.ChangeReservationRequest) (*inventoryv1.OperationResult, error) {
	return s.mutate(ctx, Request{DatasetID: r.GetDatasetId(), OrderID: r.GetOrderId(), OperationID: r.GetOperationId(), Method: "release"})
}
func (s *Server) mutate(ctx context.Context, r Request) (*inventoryv1.OperationResult, error) {
	if _, _, err := canonical(r); err != nil {
		return nil, rpcError(err)
	}
	receipt, err := s.backend.Mutate(ctx, r)
	if err != nil {
		return nil, rpcError(err)
	}
	return &inventoryv1.OperationResult{OperationId: receipt.OperationID, Outcome: inventoryv1.Outcome(inventoryv1.Outcome_value["OUTCOME_"+receipt.Outcome]), ReservationAtOperation: protoReservation(receipt.Reservation)}, nil
}
func (s *Server) GetReservation(ctx context.Context, r *inventoryv1.GetReservationRequest) (*inventoryv1.Reservation, error) {
	result, err := s.backend.GetReservation(ctx, r.GetDatasetId(), r.GetOrderId())
	if err != nil {
		return nil, rpcError(err)
	}
	return protoReservation(result), nil
}
func (s *Server) GetAvailability(ctx context.Context, r *inventoryv1.GetAvailabilityRequest) (*inventoryv1.GetAvailabilityResponse, error) {
	result, err := s.backend.GetAvailability(ctx, r.GetDatasetId(), r.GetSkus())
	if err != nil {
		return nil, rpcError(err)
	}
	response := &inventoryv1.GetAvailabilityResponse{}
	for _, item := range result {
		response.Items = append(response.Items, &inventoryv1.Availability{Sku: item.SKU, Available: item.Available, UnitPriceMinor: item.UnitPriceMinor, Currency: item.Currency})
	}
	return response, nil
}
func protoReservation(r *Reservation) *inventoryv1.Reservation {
	if r == nil {
		return nil
	}
	result := &inventoryv1.Reservation{DatasetId: r.DatasetID, OrderId: r.OrderID, State: inventoryv1.ReservationState(inventoryv1.ReservationState_value["RESERVATION_STATE_"+r.State]), Version: r.Version, TotalMinor: r.TotalMinor, Currency: r.Currency}
	for _, item := range r.Items {
		result.Items = append(result.Items, &inventoryv1.PricedItem{Sku: item.SKU, Quantity: item.Quantity, UnitPriceMinor: item.UnitPriceMinor})
	}
	return result
}
func rpcError(err error) error {
	switch {
	case errors.Is(err, ErrInvalid):
		return status.Error(codes.InvalidArgument, ErrInvalid.Error())
	case errors.Is(err, ErrKeyConflict):
		return status.Error(codes.AlreadyExists, ErrKeyConflict.Error())
	case errors.Is(err, ErrDataset):
		return status.Error(codes.FailedPrecondition, ErrDataset.Error())
	case errors.Is(err, ErrNotFound):
		return status.Error(codes.NotFound, ErrNotFound.Error())
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "inventory request canceled; outcome may be unknown")
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "inventory deadline exceeded; outcome may be unknown")
	default:
		return status.Error(codes.Unavailable, "inventory storage unavailable; retry with the same operation key")
	}
}
