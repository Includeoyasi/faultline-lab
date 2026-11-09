package orders

import (
	"context"
	"errors"

	ordersv1 "github.com/Includeoyasi/faultline-lab/api/orders/v1"
	"github.com/Includeoyasi/faultline-lab/internal/orders/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Backend interface {
	Create(context.Context, model.Request) (model.Acceptance, error)
	GetOrder(context.Context, string, string) (model.Order, error)
	GetByRequestKey(context.Context, string, string) (model.Order, error)
}
type Server struct {
	ordersv1.UnimplementedOrdersServer
	backend Backend
}

func Register(server *grpc.Server, backend Backend) {
	ordersv1.RegisterOrdersServer(server, &Server{backend: backend})
}
func (s *Server) CreateOrder(ctx context.Context, r *ordersv1.CreateOrderRequest) (*ordersv1.Acceptance, error) {
	request := model.Request{DatasetID: r.GetDatasetId(), RequestKey: r.GetRequestKey()}
	for _, item := range r.GetItems() {
		request.Items = append(request.Items, model.Item{SKU: item.GetSku(), Quantity: item.GetQuantity()})
	}
	if _, _, err := model.Canonical(request); err != nil {
		return nil, rpcError(err)
	}
	result, err := s.backend.Create(ctx, request)
	if err != nil {
		return nil, rpcError(err)
	}
	return &ordersv1.Acceptance{DatasetId: result.DatasetID, OrderId: result.OrderID, RequestKey: result.RequestKey, State: ordersv1.OrderState_ORDER_STATE_ACCEPTED}, nil
}
func (s *Server) GetOrder(ctx context.Context, r *ordersv1.GetOrderRequest) (*ordersv1.Order, error) {
	if !model.ValidIdentifier(r.GetDatasetId()) || !model.ValidIdentifier(r.GetOrderId()) {
		return nil, rpcError(model.ErrInvalid)
	}
	o, err := s.backend.GetOrder(ctx, r.GetDatasetId(), r.GetOrderId())
	if err != nil {
		return nil, rpcError(err)
	}
	return protoOrder(o), nil
}
func (s *Server) GetOrderByRequestKey(ctx context.Context, r *ordersv1.GetOrderByRequestKeyRequest) (*ordersv1.Order, error) {
	if !model.ValidIdentifier(r.GetDatasetId()) || !model.ValidIdentifier(r.GetRequestKey()) {
		return nil, rpcError(model.ErrInvalid)
	}
	o, err := s.backend.GetByRequestKey(ctx, r.GetDatasetId(), r.GetRequestKey())
	if err != nil {
		return nil, rpcError(err)
	}
	return protoOrder(o), nil
}
func protoOrder(o model.Order) *ordersv1.Order {
	result := &ordersv1.Order{DatasetId: o.DatasetID, OrderId: o.OrderID, RequestKey: o.RequestKey, State: ordersv1.OrderState(ordersv1.OrderState_value["ORDER_STATE_"+o.State]), Version: o.Version, Rejection: o.Rejection, BlockedReason: o.BlockedReason, ReserveOperationId: o.OperationID}
	for _, item := range o.Items {
		result.Items = append(result.Items, &ordersv1.Item{Sku: item.SKU, Quantity: item.Quantity})
	}
	if o.Reservation != nil {
		result.TotalMinor = o.Reservation.TotalMinor
		result.Currency = o.Reservation.Currency
		for _, item := range o.Reservation.Items {
			result.PricedItems = append(result.PricedItems, &ordersv1.PricedItem{Sku: item.SKU, Quantity: item.Quantity, UnitPriceMinor: item.UnitPriceMinor})
		}
	}
	return result
}
func rpcError(err error) error {
	switch {
	case errors.Is(err, model.ErrInvalid):
		return status.Error(codes.InvalidArgument, model.ErrInvalid.Error())
	case errors.Is(err, model.ErrDataset):
		return status.Error(codes.FailedPrecondition, model.ErrDataset.Error())
	case errors.Is(err, model.ErrKeyConflict):
		return status.Error(codes.AlreadyExists, model.ErrKeyConflict.Error())
	case errors.Is(err, model.ErrNotFound):
		return status.Error(codes.NotFound, model.ErrNotFound.Error())
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "orders request canceled; outcome may be unknown")
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "orders deadline exceeded; outcome may be unknown")
	default:
		return status.Error(codes.Unavailable, "orders storage unavailable; retry the same request key")
	}
}
