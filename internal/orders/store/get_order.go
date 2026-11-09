package store

import (
	"context"

	"github.com/Includeoyasi/faultline-lab/internal/orders/model"
)

const getOrderSQL = readOrderSQL + " WHERE o.order_id=$1"

func (s *Store) GetOrder(ctx context.Context, dataset, id string) (model.Order, error) {
	if err := s.validateLookup(dataset, id); err != nil {
		return model.Order{}, err
	}
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.read(s.pool.QueryRow(ctx, getOrderSQL, id))
}
