package inventory

import (
	"context"
	"io"
	"log/slog"

	"github.com/Includeoyasi/faultline-lab/internal/inventory/store"
	"github.com/Includeoyasi/faultline-lab/internal/service"
	"google.golang.org/grpc"
)

// Execute starts the business API by default. Health-only is an explicit
// skeleton command for dependency-free lifecycle checks, never business readiness.
func Execute(ctx context.Context, args []string, lookup func(string) (string, bool), output io.Writer) int {
	if len(args) > 0 {
		if len(args) == 1 && args[0] == "health-only" {
			return service.Execute(ctx, "inventory", nil, lookup, output)
		}
		return service.Execute(ctx, "inventory", args, lookup, output)
	}
	logger := slog.New(slog.NewJSONHandler(output, nil)).With("service", "inventory")
	// Validate listener/log/lifecycle configuration before touching the database.
	if _, err := service.LoadConfig(lookup); err != nil {
		logger.ErrorContext(ctx, "configuration rejected", "error", err.Error())
		return 1
	}
	path, _ := lookup("FL_INVENTORY_CONFIG")
	dataset, _ := lookup("FL_INVENTORY_DATASET_ID")
	cfg, err := store.LoadConfig(path, dataset)
	if err != nil {
		logger.ErrorContext(ctx, "configuration rejected", "error", err.Error())
		return 1
	}
	backend, err := store.Open(ctx, cfg)
	if err != nil {
		logger.ErrorContext(ctx, "inventory startup failed", "error", err.Error())
		return 1
	}
	defer backend.Close()
	if err = service.RunRegistered(ctx, "inventory", lookup, output, func(server *grpc.Server) { Register(server, backend) }); err != nil {
		return 1
	}
	return 0
}
