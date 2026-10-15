package service

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestHealthcheckCommand(t *testing.T) {
	t.Parallel()
	fixture := startShutdownFixture(t, time.Second)
	lookup := environment(map[string]string{"FL_HEALTH_TARGET": fixture.listener.Addr().String()})
	var output bytes.Buffer
	if code := Execute(t.Context(), "fixture", []string{"healthcheck"}, lookup, &output); code != 0 || output.Len() != 0 {
		t.Fatalf("healthy probe = %d; output %s", code, &output)
	}
	fixture.health.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	if code := Execute(t.Context(), "fixture", []string{"healthcheck"}, lookup, &output); code != 1 || !strings.Contains(output.String(), "not SERVING") {
		t.Fatalf("unhealthy probe = %d; output %s", code, &output)
	}
}

func TestHealthcheckRejectsInvalidConfiguration(t *testing.T) {
	t.Parallel()
	for _, target := range []string{"", ":50051", "0.0.0.0:50051", "[::]:50051", "secret://user:password@host:80"} {
		t.Run(target, func(t *testing.T) {
			var output bytes.Buffer
			code := Execute(t.Context(), "orders", []string{"healthcheck"}, environment(map[string]string{"FL_HEALTH_TARGET": target}), &output)
			if code != 1 || !strings.Contains(output.String(), "FL_HEALTH_TARGET is required") || (target != "" && strings.Contains(output.String(), target)) {
				t.Fatalf("probe = %d; output %s", code, &output)
			}
		})
	}
	var output bytes.Buffer
	if code := Execute(t.Context(), "orders", []string{"secret-argument"}, environment(nil), &output); code != 2 || strings.Contains(output.String(), "secret-argument") {
		t.Fatalf("unknown command = %d; output %s", code, &output)
	}
}

func TestHealthcheckHonorsCancellation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	var output bytes.Buffer
	lookup := environment(map[string]string{"FL_HEALTH_TARGET": "localhost:1"})
	if code := Execute(ctx, "orders", []string{"healthcheck"}, lookup, &output); code != 1 || !strings.Contains(output.String(), "health check failed") {
		t.Fatalf("cancelled probe = %d; output %s", code, &output)
	}
}
