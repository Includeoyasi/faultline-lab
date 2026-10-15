package service

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"runtime/debug"
	"testing"
	"time"
)

func TestJSONLoggingAndLevel(t *testing.T) {
	t.Parallel()
	var output bytes.Buffer
	logger := newLogger(&output, "orders", slog.LevelWarn)
	logger.Info("hidden")
	if output.Len() != 0 {
		t.Fatal("INFO was not filtered at WARN")
	}
	logger.Error("bind failed", "error", "address already in use")
	var record map[string]any
	if err := json.Unmarshal(output.Bytes(), &record); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if record["service"] != "orders" || record["level"] != "ERROR" || record["msg"] != "bind failed" || record["error"] == "" {
		t.Fatalf("unexpected record: %v", record)
	}
	if _, err := time.Parse(time.RFC3339Nano, record["time"].(string)); err != nil {
		t.Fatalf("invalid log timestamp: %v", err)
	}
}

func TestBuildAttributes(t *testing.T) {
	t.Parallel()
	for _, available := range []bool{false, true} {
		t.Run(map[bool]string{false: "unknown", true: "available"}[available], func(t *testing.T) {
			info := &debug.BuildInfo{GoVersion: "go1.test", Main: debug.Module{Version: "v0.test"}, Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "actual-revision"}, {Key: "vcs.modified", Value: "true"},
			}}
			var output bytes.Buffer
			newLogger(&output, "inventory", slog.LevelInfo).LogAttrs(context.Background(), slog.LevelInfo, "service started", buildAttributes(info, available)...)
			var record map[string]any
			if err := json.Unmarshal(output.Bytes(), &record); err != nil {
				t.Fatal(err)
			}
			for key, value := range map[string]string{"version": "v0.test", "revision": "actual-revision", "vcs_modified": "true", "go_version": "go1.test"} {
				if !available {
					value = "unknown"
				}
				if record[key] != value {
					t.Errorf("%s = %v, want %s", key, record[key], value)
				}
			}
		})
	}
}
