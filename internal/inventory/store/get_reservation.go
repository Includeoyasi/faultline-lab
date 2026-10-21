package store

import (
	"context"

	"github.com/Includeoyasi/faultline-lab/internal/inventory/model"
)

func (s *Store) GetReservation(ctx context.Context, dataset, orderID string) (*model.Reservation, error) {
	if dataset != s.dataset {
		return nil, model.ErrDataset
	}
	if !model.ValidIdentifier(orderID) {
		return nil, model.ErrInvalid
	}
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.readReservation(ctx, s.pool, orderID)
}
