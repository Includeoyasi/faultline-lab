package service

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	testpb "google.golang.org/grpc/interop/grpc_testing"
	"google.golang.org/grpc/status"
)

type closingListener struct {
	net.Listener
	closed chan struct{}
	once   sync.Once
}

func (listener *closingListener) Close() error {
	err := listener.Listener.Close()
	listener.once.Do(func() { close(listener.closed) })
	return err
}

type controlledRPC struct {
	testpb.UnimplementedTestServiceServer
	entered chan struct{}
	release chan struct{}
	exited  chan struct{}
}

func (rpc *controlledRPC) EmptyCall(context.Context, *testpb.Empty) (*testpb.Empty, error) {
	close(rpc.entered)
	defer close(rpc.exited)
	// Deliberately ignores cancellation to prove that shutdown returns on its
	// deadline. The test owns release and always lets this handler exit afterward.
	<-rpc.release
	return &testpb.Empty{}, nil
}

type shutdownFixture struct {
	listener      *closingListener
	server        *grpc.Server
	health        *health.Server
	work          *controlledRPC
	conn          *grpc.ClientConn
	cancel        context.CancelFunc
	stopping      chan struct{}
	continueDrain func()
	releaseWork   func()
	done          chan struct{}
	runErr        error // Read after done closes.
	logs          bytes.Buffer
}

func startShutdownFixture(t *testing.T, timeout time.Duration) *shutdownFixture {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	gate := make(chan struct{})
	fixture := &shutdownFixture{
		listener: &closingListener{Listener: listener, closed: make(chan struct{})},
		server:   grpc.NewServer(), health: health.NewServer(), cancel: cancel,
		stopping: make(chan struct{}), done: make(chan struct{}),
		work: &controlledRPC{entered: make(chan struct{}), release: make(chan struct{}), exited: make(chan struct{})},
	}
	fixture.continueDrain = sync.OnceFunc(func() { close(gate) })
	fixture.releaseWork = sync.OnceFunc(func() { close(fixture.work.release) })
	healthpb.RegisterHealthServer(fixture.server, fixture.health)
	testpb.RegisterTestServiceServer(fixture.server, fixture.work)
	writer := writerFunc(func(data []byte) (int, error) {
		length, err := fixture.logs.Write(data)
		if strings.Contains(string(data), `"msg":"service stopping"`) {
			close(fixture.stopping)
			// Observe real health over TCP before GracefulStop closes the listener.
			<-gate
		}
		return length, err
	})
	go func() {
		fixture.runErr = serveGRPC(ctx, fixture.listener, fixture.server, fixture.health, timeout, newLogger(writer, "fixture", slog.LevelInfo))
		close(fixture.done)
	}()
	t.Cleanup(func() {
		cancel()
		fixture.continueDrain()
		fixture.releaseWork()
		await(t, fixture.done, "server cleanup")
		fixture.server.Stop()
		if fixture.conn != nil {
			if err := fixture.conn.Close(); err != nil {
				t.Error(err)
			}
		}
	})
	fixture.conn, err = grpc.NewClient("passthrough:///"+listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	assertHealth(t, fixture.conn, healthpb.HealthCheckResponse_SERVING)
	return fixture
}

func await(t *testing.T, signal <-chan struct{}, operation string) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for %s", operation)
	}
}

func assertHealth(t *testing.T, conn *grpc.ClientConn, want healthpb.HealthCheckResponse_ServingStatus) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()
	response, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{}, grpc.WaitForReady(true))
	if err != nil || response.GetStatus() != want {
		t.Fatalf("health = %v, %v; want %v", response, err, want)
	}
}

func beginWork(t *testing.T, fixture *shutdownFixture) <-chan error {
	t.Helper()
	result := make(chan error, 1)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)
	go func() {
		_, err := testpb.NewTestServiceClient(fixture.conn).EmptyCall(ctx, &testpb.Empty{})
		result <- err
	}()
	await(t, fixture.work.entered, "active RPC")
	return result
}

func beginDrain(t *testing.T, fixture *shutdownFixture) {
	t.Helper()
	fixture.cancel()
	await(t, fixture.stopping, "readiness removal")
	assertHealth(t, fixture.conn, healthpb.HealthCheckResponse_NOT_SERVING)
	fixture.continueDrain()
	await(t, fixture.listener.closed, "listener closure")
	connection, err := net.DialTimeout("tcp", fixture.listener.Addr().String(), time.Second)
	if err == nil {
		_ = connection.Close()
		t.Fatal("listener accepted a new connection during drain")
	}
}

func TestGracefulShutdownFinishesActiveRPC(t *testing.T) {
	t.Parallel()
	fixture := startShutdownFixture(t, 3*time.Second)
	result := beginWork(t, fixture)
	beginDrain(t, fixture)
	select {
	case err := <-result:
		t.Fatalf("active RPC ended before work was released: %v", err)
	default:
	}
	fixture.releaseWork()
	select {
	case err := <-result:
		if err != nil {
			t.Fatalf("active RPC was interrupted: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("released RPC did not complete")
	}
	await(t, fixture.done, "graceful shutdown")
	if fixture.runErr != nil || !strings.Contains(fixture.logs.String(), `"stop_mode":"graceful"`) {
		t.Fatalf("shutdown = %v; logs: %s", fixture.runErr, &fixture.logs)
	}
}

func TestShutdownDeadlineInterruptsBlockedRPC(t *testing.T) {
	t.Parallel()
	const budget = 200 * time.Millisecond
	fixture := startShutdownFixture(t, budget)
	result := beginWork(t, fixture)
	started := time.Now()
	beginDrain(t, fixture)
	await(t, fixture.done, "bounded forced shutdown")
	if !errors.Is(fixture.runErr, ErrShutdownTimeout) || time.Since(started) < budget || !strings.Contains(fixture.logs.String(), `"stop_mode":"forced"`) {
		t.Fatalf("shutdown = %v after %v; logs: %s", fixture.runErr, time.Since(started), &fixture.logs)
	}
	select {
	case err := <-result:
		if status.Code(err) != codes.Unavailable {
			t.Fatalf("forced RPC = %v, want Unavailable", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("client did not observe forced interruption")
	}
	// Release the intentionally noncooperative handler so the test also joins
	// gRPC internals. Real cmd entrypoints exit with code 1 on this fatal path.
	fixture.releaseWork()
	await(t, fixture.work.exited, "blocked handler cleanup")
}
