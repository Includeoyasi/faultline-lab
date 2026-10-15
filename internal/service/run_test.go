package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestRunRejectsConfigurationBeforeListening(t *testing.T) {
	t.Parallel()
	var output bytes.Buffer
	err := Run(context.Background(), "payments", environment(map[string]string{"FL_GRPC_ADDR": "localhost:50053", "FL_LOG_LEVEL": "super-secret-value"}), &output)
	if err == nil || strings.Contains(output.String(), "super-secret-value") {
		t.Fatalf("expected sanitized configuration failure, got %v / %s", err, &output)
	}
	var record map[string]any
	if err := json.Unmarshal(output.Bytes(), &record); err != nil {
		t.Fatal(err)
	}
	if record["service"] != "payments" || record["level"] != "ERROR" || record["msg"] != "configuration rejected" {
		t.Fatalf("unexpected diagnostic: %v", record)
	}
}

func TestRunRejectsOccupiedAddress(t *testing.T) {
	t.Parallel()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := listener.Close(); err != nil {
			t.Error(err)
		}
	})
	var output bytes.Buffer
	err = Run(context.Background(), "orders", environment(map[string]string{"FL_GRPC_ADDR": listener.Addr().String()}), &output)
	var operationError *net.OpError
	if !errors.As(err, &operationError) || operationError.Op != "listen" {
		t.Fatalf("expected wrapped listen error, got %v", err)
	}
	var record map[string]any
	if err := json.Unmarshal(output.Bytes(), &record); err != nil {
		t.Fatal(err)
	}
	if record["service"] != "orders" || record["msg"] != "service failed" || record["level"] != "ERROR" {
		t.Fatalf("unexpected bind diagnostic: %v", record)
	}
}

func TestServeWithCancelledContextDoesNotAnnounceReadiness(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var output bytes.Buffer
	err := serve(ctx, Config{GRPCAddress: "127.0.0.1:0"}, newLogger(&output, "orders", slog.LevelInfo))
	if !errors.Is(err, context.Canceled) || output.Len() != 0 {
		t.Fatalf("expected cancellation without startup, got %v / %s", err, &output)
	}
}

type recordWriter chan []byte

func (records recordWriter) Write(data []byte) (int, error) {
	records <- append([]byte(nil), data...)
	return len(data), nil
}

type writerFunc func([]byte) (int, error)

func (write writerFunc) Write(data []byte) (int, error) { return write(data) }

func TestServeCancellationDuringStartup(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	var output bytes.Buffer
	writer := writerFunc(func(data []byte) (int, error) {
		// Cancel before the Serve goroutine is started. Stop may win that race;
		// grpc.ErrServerStopped is then an expected shutdown result.
		cancel()
		return output.Write(data)
	})
	if err := serve(ctx, Config{GRPCAddress: "127.0.0.1:0"}, newLogger(writer, "orders", slog.LevelInfo)); err != nil {
		t.Fatalf("cancellation during startup failed: %v", err)
	}
	if !strings.Contains(output.String(), `"msg":"service stopped"`) {
		t.Fatalf("missing final log: %s", &output)
	}
}

func TestServeHealthAndCancellation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	records := make(recordWriter, 4)
	result := make(chan error, 1)
	go func() {
		// Port zero is only a test fixture; public configuration rejects it.
		result <- serve(ctx, Config{GRPCAddress: "127.0.0.1:0"}, newLogger(records, "orders", slog.LevelInfo))
	}()
	t.Cleanup(func() {
		cancel()
		select {
		case err := <-result:
			if err != nil {
				t.Errorf("serve: %v", err)
			}
		case <-time.After(10 * time.Second):
			t.Error("server did not exit after cancellation")
		}
	})
	var started map[string]any
	select {
	case record := <-records:
		if err := json.Unmarshal(record, &started); err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal("missing startup record")
	}
	address, ok := started["grpc_address"].(string)
	if !ok || started["msg"] != "service started" {
		t.Fatalf("invalid startup record: %v", started)
	}
	conn, err := grpc.NewClient("passthrough:///"+address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()
	response, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{}, grpc.WaitForReady(true))
	if err != nil || response.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("health = %v, %v", response, err)
	}
	cancel()
	select {
	case record := <-records:
		var stopped map[string]any
		if err := json.Unmarshal(record, &stopped); err != nil {
			t.Fatal(err)
		}
		if stopped["msg"] != "service stopped" || stopped["service"] != "orders" || stopped["stop_mode"] != "immediate" {
			t.Fatalf("invalid stop record: %v", stopped)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("missing stop record")
	}
	// Cleanup owns the single receive from result, so it always joins the server.
}
