package service

import (
	"log/slog"
	"strings"
	"testing"
	"time"
)

func environment(values map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) { value, ok := values[key]; return value, ok }
}

func TestLoadConfig(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		env  map[string]string
		want Config
	}{
		{"non-network defaults", map[string]string{"FL_GRPC_ADDR": "localhost:50051"}, Config{"localhost:50051", slog.LevelInfo, 10 * time.Second}},
		{"overrides", map[string]string{"FL_GRPC_ADDR": "0.0.0.0:8080", "FL_LOG_LEVEL": "debug", "FL_SHUTDOWN_TIMEOUT": "250ms"}, Config{"0.0.0.0:8080", slog.LevelDebug, 250 * time.Millisecond}},
		{"ipv6", map[string]string{"FL_GRPC_ADDR": "[::1]:8080", "FL_LOG_LEVEL": "WARN"}, Config{"[::1]:8080", slog.LevelWarn, 10 * time.Second}},
		{"dns", map[string]string{"FL_GRPC_ADDR": "orders.local:50051", "FL_LOG_LEVEL": "ERROR"}, Config{"orders.local:50051", slog.LevelError, 10 * time.Second}},
		{"wildcard", map[string]string{"FL_GRPC_ADDR": ":50051"}, Config{":50051", slog.LevelInfo, 10 * time.Second}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := LoadConfig(environment(tt.env))
			if err != nil || got != tt.want {
				t.Fatalf("LoadConfig() = %+v, %v; want %+v", got, err, tt.want)
			}
		})
	}
}

func TestLoadConfigRejectsInvalidValues(t *testing.T) {
	t.Parallel()
	inputs := map[string][]string{
		"FL_GRPC_ADDR":        {"", "localhost", "::1:8080", "localhost:http", ":0", ":65536", ":-1", ":+80", "bad host:80", "-host:80", "host..local:80", "host_:80", "secret://user:password@host:80"},
		"FL_LOG_LEVEL":        {"", "TRACE", "INFO+1", "super-secret-value"},
		"FL_SHUTDOWN_TIMEOUT": {"", "0s", "-1s", "10", "999999999999999999999h", "super-secret-value"},
	}
	for key, values := range inputs {
		for _, value := range values {
			t.Run(key+"/"+value, func(t *testing.T) {
				t.Parallel()
				values := map[string]string{"FL_GRPC_ADDR": "localhost:50051"}
				values[key] = value
				_, err := LoadConfig(environment(values))
				if err == nil || !strings.Contains(err.Error(), key) {
					t.Fatalf("expected diagnostic naming %s, got %v", key, err)
				}
				if value != "" && strings.Contains(err.Error(), value) {
					t.Fatalf("diagnostic echoed the supplied value: %v", err)
				}
			})
		}
	}
	if _, err := LoadConfig(environment(nil)); err == nil || !strings.Contains(err.Error(), "FL_GRPC_ADDR is required") {
		t.Fatalf("missing listen address must fail closed, got %v", err)
	}
}
