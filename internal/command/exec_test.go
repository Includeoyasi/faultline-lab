package command

import (
	"context"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCommandHelper(t *testing.T) {
	if os.Getenv("FAULTLINE_EXEC_HELPER") != "1" {
		return
	}
	if os.Getenv("FAULTLINE_EXEC_CHILD") == "1" {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			os.Exit(3)
		}
		if err := os.WriteFile(os.Getenv("FAULTLINE_EXEC_READY"), []byte(listener.Addr().String()), 0o600); err != nil {
			os.Exit(3)
		}
		connection, err := listener.Accept()
		if err != nil {
			os.Exit(4)
		}
		_, _ = io.Copy(io.Discard, connection)
		_ = connection.Close()
		_ = listener.Close()
		os.Exit(0)
	}
	child := exec.Command(os.Args[0], "-test.run=^TestCommandHelper$")
	child.Env = append(os.Environ(), "FAULTLINE_EXEC_CHILD=1")
	child.Stdout, child.Stderr = os.Stdout, os.Stderr
	if os.Getenv("FAULTLINE_EXEC_PARENT_EXITS") == "1" {
		child.Stdout, child.Stderr = nil, nil
	}
	if err := child.Start(); err != nil {
		os.Exit(2)
	}
	if os.Getenv("FAULTLINE_EXEC_PARENT_EXITS") == "1" {
		for {
			if _, err := os.Stat(os.Getenv("FAULTLINE_EXEC_RELEASE")); err == nil {
				os.Exit(0)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	if err := child.Wait(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func TestExecutorCancelsPluginTree(t *testing.T) {
	ready := filepath.Join(t.TempDir(), "ready")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	executor := processExecutor{directory: t.TempDir(), env: append(os.Environ(), "FAULTLINE_EXEC_HELPER=1", "FAULTLINE_EXEC_READY="+ready)}
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
	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("cancelled command succeeded")
		}
	case <-time.After(7 * time.Second):
		t.Fatal("command tree did not settle")
	}
	// A living plugin retains this socket even after its parent was killed.
	if err = connection.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatal(err)
	}
	var value [1]byte
	_, err = connection.Read(value[:])
	if networkError, ok := err.(net.Error); err == nil || (ok && networkError.Timeout()) {
		t.Fatalf("plugin survived command cancellation: %v", err)
	}
}

func TestExecutorMissingProgram(t *testing.T) {
	executor := processExecutor{directory: t.TempDir(), env: os.Environ()}
	output, err := executor.execute(context.Background(), "faultline-nonexistent-test-command", nil)
	if err == nil || output.exitCode != -1 {
		t.Fatalf("output=%+v err=%v", output, err)
	}
}

func TestOutputCaptureIsBounded(t *testing.T) {
	buffer := &boundedBuffer{}
	data := make([]byte, outputLimit+123)
	if count, err := buffer.Write(data); err != nil || count != len(data) || buffer.Len() != outputLimit || !buffer.truncated {
		t.Fatalf("unbounded output: n=%d err=%v len=%d", count, err, buffer.Len())
	}
	buffer = &boundedBuffer{}
	reader := struct{ io.Reader }{strings.NewReader(strings.Repeat("x", outputLimit+123))}
	if count, err := io.Copy(buffer, reader); err != nil || count != outputLimit+123 || buffer.Len() != outputLimit || !buffer.truncated {
		t.Fatalf("io.Copy bypassed limit: n=%d err=%v len=%d", count, err, buffer.Len())
	}
}
