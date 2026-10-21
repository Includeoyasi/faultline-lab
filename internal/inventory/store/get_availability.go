package store

import (
	"context"

	"github.com/Includeoyasi/faultline-lab/internal/inventory/model"
)

const (
	readAvailabilitySQL = "SELECT s.sku,s.available,p.price_minor,p.currency FROM inventory.stock s JOIN inventory.products p USING(sku) WHERE s.sku=ANY($1) ORDER BY s.sku"
)

func (s *Store) GetAvailability(ctx context.Context, dataset string, skus []string) ([]model.Availability, error) {
	if dataset != s.dataset {
		return nil, model.ErrDataset
	}
	if len(skus) == 0 || len(skus) > 100 {
		return nil, model.ErrInvalid
	}
	seen := map[string]bool{}
	for _, sku := range skus {
		if !model.ValidSKU(sku) || seen[sku] {
			return nil, model.ErrInvalid
		}
		seen[sku] = true
	}
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	rows, err := s.pool.Query(ctx, readAvailabilitySQL, skus)
	if err != nil {
		return nil, storageError(err)
	}
	defer rows.Close()
	result := []model.Availability{}
	for rows.Next() {
		var item model.Availability
		if err = rows.Scan(&item.SKU, &item.Available, &item.UnitPriceMinor, &item.Currency); err != nil {
			return nil, storageError(err)
		}
		result = append(result, item)
	}
	if err = rows.Err(); err != nil {
		return nil, storageError(err)
	}
	if len(result) != len(skus) {
		return nil, model.ErrNotFound
	}
	return result, nil
}
