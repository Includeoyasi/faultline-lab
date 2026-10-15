package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// ErrShutdownTimeout indicates that active RPCs required a forced shutdown.
var ErrShutdownTimeout = errors.New("shutdown timeout exceeded; active RPCs were interrupted")

// Run validates configuration, serves gRPC health and drains active RPCs on
// cancellation. It returns ErrShutdownTimeout if the drain budget expires.
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
	return serveGRPC(ctx, listener, server, healthServer, cfg.ShutdownTimeout, logger)
}

// serveGRPC owns the server, listener and all shutdown goroutines. The explicit
// server parameter also lets lifecycle tests register a controlled RPC handler.
func serveGRPC(ctx context.Context, listener net.Listener, server *grpc.Server, healthServer *health.Server, timeout time.Duration, logger *slog.Logger) error {
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	info, ok := debug.ReadBuildInfo()
	attrs := append(buildAttributes(info, ok), slog.String("grpc_address", listener.Addr().String()))
	logger.LogAttrs(ctx, slog.LevelInfo, "service started", attrs...)
	serveResult := make(chan error, 1)
	go func() { serveResult <- server.Serve(listener) }()
	select {
	case err := <-serveResult:
		server.Stop()
		if err == nil {
			err = errors.New("server stopped unexpectedly")
		}
		return fmt.Errorf("serve gRPC: %w", err)
	case <-ctx.Done():
		healthServer.Shutdown()
		logger.InfoContext(ctx, "service stopping", "health", "NOT_SERVING", "shutdown_timeout", timeout.String())
		drained := make(chan struct{})
		go func() { server.GracefulStop(); close(drained) }()
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		select {
		case <-drained:
		case <-timer.C:
			// Give a completed drain precedence if completion and deadline race.
			select {
			case <-drained:
			default:
				// GracefulStop can wait for a handler that ignores cancellation.
				// Do not join it on this fatal path: Run returns an error and the
				// entrypoint exits the process, terminating remaining RPC work.
				go server.Stop()
				logger.WarnContext(ctx, "service stopped", "stop_mode", "forced")
				return ErrShutdownTimeout
			}
		}
		if err := <-serveResult; err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return fmt.Errorf("stop gRPC: %w", err)
		}
		logger.InfoContext(ctx, "service stopped", "stop_mode", "graceful")
		return nil
	}
}
