package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Run validates configuration, serves standard gRPC health and blocks until
// cancellation or a serving error. Cancellation currently stops RPCs immediately;
// draining active work with a deadline is the separate P01-05 task.
func Run(ctx context.Context, name string, lookup func(string) (string, bool), output io.Writer) error {
	logger := newLogger(output, name, slog.LevelInfo)
	cfg, err := LoadConfig(lookup)
	if err != nil {
		logger.ErrorContext(ctx, "configuration rejected", "error", err.Error())
		return err
	}
	logger = newLogger(output, name, cfg.LogLevel)
	if err := serve(ctx, cfg, logger); err != nil {
		logger.ErrorContext(ctx, "service failed", "error", err.Error())
		return err
	}
	return nil
}

func serve(ctx context.Context, cfg Config, logger *slog.Logger) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	listener, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		return fmt.Errorf("listen for gRPC: %w", err)
	}
	// grpc.Server.Serve owns and closes the listener, including on Stop.
	server := grpc.NewServer()
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	healthpb.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	info, ok := debug.ReadBuildInfo()
	attrs := append(buildAttributes(info, ok), slog.String("grpc_address", listener.Addr().String()))
	logger.LogAttrs(ctx, slog.LevelInfo, "service started", attrs...)
	serveResult := make(chan error, 1)
	go func() { serveResult <- server.Serve(listener) }()
	select {
	case err := <-serveResult:
		server.Stop()
		return fmt.Errorf("serve gRPC: %w", err)
	case <-ctx.Done():
		healthServer.Shutdown()
		server.Stop()
		if err := <-serveResult; err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return fmt.Errorf("stop gRPC: %w", err)
		}
		logger.InfoContext(ctx, "service stopped", "stop_mode", "immediate")
		return nil
	}
}
