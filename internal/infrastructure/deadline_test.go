package infrastructure

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"
	"time"
)

// Accept connections but never reply, so refusal/DNS does not mask a missing deadline.
func blackhole(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	var connections []net.Conn
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			mu.Lock()
			connections = append(connections, conn)
			mu.Unlock()
		}
	}()
	t.Cleanup(func() {
		_ = listener.Close()
		<-done
		mu.Lock()
		defer mu.Unlock()
		for _, conn := range connections {
			_ = conn.Close()
		}
	})
	return listener.Addr().String()
}

func TestConnectionDeadlines(t *testing.T) {
	for _, dependency := range []string{"redis", "kafka"} {
		t.Run(dependency, func(t *testing.T) {
			cfg := testConfig(t)
			endpoint := blackhole(t)
			ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
			defer cancel()
			started := time.Now()
			var err error
			if dependency == "redis" {
				cfg.RedisAddress = endpoint
				client := redisClient(cfg, "lab", "unusable")
				err = client.Ping(ctx).Err()
				_ = client.Close()
			} else {
				cfg.KafkaBrokers = []string{endpoint}
				client, e := kafkaClient(ctx, cfg)
				if e != nil {
					t.Fatal(e)
				}
				err = client.Ping(ctx)
				client.Close()
			}
			finished := time.Now()
			deadline, _ := ctx.Deadline()
			var timeout net.Error
			deadlineError := errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) || (errors.As(err, &timeout) && timeout.Timeout())
			// Socket deadlines and the context timer can fire independently.
			// The operation must time out at its deadline, even if ctx.Err is not set yet.
			if err == nil || !deadlineError || finished.Before(deadline) || finished.Sub(started) > time.Second {
				t.Fatalf("dependency deadline not obeyed: err=%v (%T), ctx=%v, elapsed=%v", err, err, ctx.Err(), finished.Sub(started))
			}
		})
	}
}
