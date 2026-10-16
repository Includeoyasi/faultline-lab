//go:build !windows

package command

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExecutorStopsPluginAfterSuccessfulParent(t *testing.T) {
	directory := t.TempDir()
	ready, release := filepath.Join(directory, "ready"), filepath.Join(directory, "release")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	executor := processExecutor{directory: directory, env: append(os.Environ(),
		"FAULTLINE_EXEC_HELPER=1", "FAULTLINE_EXEC_PARENT_EXITS=1",
		"FAULTLINE_EXEC_READY="+ready, "FAULTLINE_EXEC_RELEASE="+release)}
	done := make(chan error, 1)
	go func() {
		_, err := executor.execute(ctx, os.Args[0], []string{"-test.run=^TestCommandHelper$"})
		done <- err
	}()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	var address []byte
	for {
		if data, err := os.ReadFile(ready); err == nil && len(data) > 0 {
			address = data
			break
		}
		select {
		case err := <-done:
			t.Fatalf("helper exited before child readiness: %v", err)
		case <-ctx.Done():
			t.Fatal("child readiness timed out")
		case <-ticker.C:
		}
	}
	connection, err := net.DialTimeout("tcp", string(address), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = connection.Close() }()
	if err = os.WriteFile(release, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	select {
	case err = <-done:
		if err != nil {
			t.Fatalf("successful parent: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("executor retained an orphan plugin")
	}
	if err = connection.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatal(err)
	}
	var value [1]byte
	_, err = connection.Read(value[:])
	if networkError, ok := err.(net.Error); err == nil || (ok && networkError.Timeout()) {
		t.Fatalf("plugin survived successful parent: %v", err)
	}
}
