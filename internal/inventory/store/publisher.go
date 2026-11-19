package store

import (
	"context"
	"log/slog"

	"github.com/Includeoyasi/faultline-lab/internal/publishing"
)

func (s *Store) Publish(ctx context.Context, broker string, logger *slog.Logger) {
	p := publishing.Publisher{Pool: s.pool, Owner: "inventory", Dataset: s.dataset, Broker: broker}
	p.Run(ctx, logger)
}
