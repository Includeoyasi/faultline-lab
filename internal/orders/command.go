// Package orders serves admission/read APIs and a bounded durable reserve worker.
package orders

import (
	"context"
	"io"
	"log/slog"
	"time"

	inventoryv1 "github.com/Includeoyasi/faultline-lab/api/inventory/v1"
	"github.com/Includeoyasi/faultline-lab/internal/orders/store"
	"github.com/Includeoyasi/faultline-lab/internal/orders/worker"
	"github.com/Includeoyasi/faultline-lab/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Execute(ctx context.Context, args []string, lookup func(string) (string, bool), output io.Writer) int {
	if len(args) > 0 {
		if len(args) == 1 && args[0] == "health-only" {
			return service.Execute(ctx, "orders", nil, lookup, output)
		}
		return service.Execute(ctx, "orders", args, lookup, output)
	}
	logger := slog.New(slog.NewJSONHandler(output, nil)).With("service", "orders")
	if _, err := service.LoadConfig(lookup); err != nil {
		logger.ErrorContext(ctx, "configuration rejected", "error", err.Error())
		return 1
	}
	path, _ := lookup("FL_ORDERS_CONFIG")
	dataset, _ := lookup("FL_ORDERS_DATASET_ID")
	cfg, err := LoadConfig(path, dataset)
	if err != nil {
		logger.ErrorContext(ctx, "configuration rejected", "error", err.Error())
		return 1
	}
	backend, err := store.Open(ctx, cfg.Database)
	if err != nil {
		logger.ErrorContext(ctx, "orders startup failed", "error", err.Error())
		return 1
	}
	defer backend.Close()
	conn, err := grpc.NewClient("passthrough:///"+cfg.InventoryTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorContext(ctx, "cannot configure Inventory client")
		return 1
	}
	defer func() { _ = conn.Close() }()
	duration := func(value string) time.Duration { v, _ := time.ParseDuration(value); return v }
	processor := worker.Processor{Jobs: backend, Inventory: inventoryv1.NewInventoryClient(conn), Lease: duration(cfg.Lease), RPCTimeout: duration(cfg.RPCTimeout), RetryMin: duration(cfg.RetryMin), RetryMax: duration(cfg.RetryMax), Logger: logger}
	workersCtx, stop := context.WithCancel(ctx)
	defer stop()
	done := make(chan struct{})
	go func() { processor.Run(workersCtx, cfg.WorkerCount, duration(cfg.PollInterval)); close(done) }()
	err = service.RunRegistered(ctx, "orders", lookup, output, func(server *grpc.Server) { Register(server, backend) })
	stop()
	<-done
	if err != nil {
		return 1
	}
	return 0
}
