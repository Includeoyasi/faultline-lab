package store

import (
	"context"

	"github.com/Includeoyasi/faultline-lab/internal/orders/model"
)

const getByRequestKeySQL = readOrderSQL + " WHERE o.request_key=$1"

func (s *Store) GetByRequestKey(ctx context.Context, dataset, key string) (model.Order, error) {
	if err := s.validateLookup(dataset, key); err != nil {
		return model.Order{}, err
	}
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.read(s.pool.QueryRow(ctx, getByRequestKeySQL, key))
}
