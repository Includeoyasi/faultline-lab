package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Execute runs the service or a bounded local health probe. It returns the exit
// code; only cmd entrypoints call os.Exit, including after a forced shutdown.
func Execute(ctx context.Context, name string, args []string, lookup func(string) (string, bool), output io.Writer) int {
	if len(args) == 0 {
		if err := Run(ctx, name, lookup, output); err != nil {
			return 1
		}
		return 0
	}
	logger := newLogger(output, name, slog.LevelInfo)
	if len(args) != 1 || args[0] != "healthcheck" {
		logger.ErrorContext(ctx, "invalid command", "usage", "service [healthcheck]")
		return 2
	}
	if err := probe(ctx, lookup); err != nil {
		logger.ErrorContext(ctx, "health check failed", "error", err.Error())
		return 1
	}
	return 0
}

func probe(ctx context.Context, lookup func(string) (string, bool)) error {
	target, _ := lookup("FL_HEALTH_TARGET")
	host, _, err := net.SplitHostPort(target)
	if err != nil || !validAddress(target) || host == "" || net.ParseIP(host).IsUnspecified() {
		return errors.New("FL_HEALTH_TARGET is required and must contain a concrete host and port")
	}
	// This command only probes the local plaintext lab. TLS for a deployed
	// business system will be an explicit, separately reviewed configuration.
	conn, err := grpc.NewClient("passthrough:///"+target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.New("create health check client")
	}
	defer func() { _ = conn.Close() }() // Sole owner; a second Close is the only error case.
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	response, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{}, grpc.WaitForReady(true))
	if err != nil {
		return errors.New("gRPC health request failed or exceeded its 2s deadline")
	}
	if response.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		return errors.New("gRPC health is not SERVING")
	}
	return nil
}
