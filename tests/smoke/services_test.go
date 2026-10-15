// Package smoke tests the actual service binaries over TCP, without databases.
package smoke

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type logCapture struct {
	mu      sync.Mutex
	pending string
	records []map[string]any
	err     error
	started chan string
}

func (capture *logCapture) Write(data []byte) (int, error) {
	capture.mu.Lock()
	defer capture.mu.Unlock()
	capture.pending += string(data)
	for {
		line, rest, found := strings.Cut(capture.pending, "\n")
		if !found {
			break
		}
		capture.pending = rest
		var record map[string]any
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			capture.err = fmt.Errorf("non-JSON process output: %w", err)
			continue
		}
		capture.records = append(capture.records, record)
		if record["msg"] == "service started" {
			address, _ := record["grpc_address"].(string)
			select {
			case capture.started <- address:
			default:
			}
		}
	}
	return len(data), nil
}

func (capture *logCapture) check(t *testing.T, service string) []map[string]any {
	t.Helper()
	capture.mu.Lock()
	defer capture.mu.Unlock()
	if capture.err != nil || capture.pending != "" {
		t.Fatalf("invalid process logs: %v; unfinished line %q", capture.err, capture.pending)
	}
	for _, record := range capture.records {
		if record["service"] != service || record["level"] == "" || record["msg"] == "" {
			t.Fatalf("missing common log fields: %v", record)
		}
		timestamp, ok := record["time"].(string)
		if !ok {
			t.Fatalf("missing timestamp: %v", record)
		}
		if _, err := time.Parse(time.RFC3339Nano, timestamp); err != nil {
			t.Fatalf("invalid timestamp: %v", err)
		}
	}
	return append([]map[string]any(nil), capture.records...)
}

type process struct {
	cmd     *exec.Cmd
	logs    *logCapture
	done    chan struct{}
	waitErr error // Read only after done closes.
	address string
}

func childEnvironment(overrides map[string]string) []string {
	var result []string
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if !strings.HasPrefix(strings.ToUpper(key), "FL_") {
			result = append(result, entry)
		}
	}
	for key, value := range overrides {
		result = append(result, key+"="+value)
	}
	return result
}

func freeAddress(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	return address
}

func start(t *testing.T, binary string) *process {
	t.Helper()
	p := &process{cmd: exec.Command(binary), logs: &logCapture{started: make(chan string, 1)}, done: make(chan struct{})}
	p.cmd.Env = childEnvironment(map[string]string{"FL_GRPC_ADDR": freeAddress(t)})
	p.cmd.Stdout, p.cmd.Stderr = p.logs, p.logs
	if err := p.cmd.Start(); err != nil {
		t.Fatal(err)
	}
	go func() { p.waitErr = p.cmd.Wait(); close(p.done) }()
	t.Cleanup(func() {
		select {
		case <-p.done:
			return
		default:
		}
		if err := p.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			t.Error(err)
		}
		select {
		case <-p.done:
		case <-time.After(10 * time.Second):
			t.Error("child did not exit after cleanup")
		}
	})
	select {
	case p.address = <-p.logs.started:
	case <-p.done:
		t.Fatalf("service exited before startup: %v; logs: %v", p.waitErr, p.logs.check(t, filepath.Base(binary)))
	case <-time.After(10 * time.Second):
		t.Fatal("service did not announce startup")
	}
	return p
}

func checkHealth(t *testing.T, address, service string) {
	t.Helper()
	conn, err := grpc.NewClient("passthrough:///"+address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	response, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{Service: service}, grpc.WaitForReady(true))
	if service == "unknown-business-api" {
		if status.Code(err) != codes.NotFound {
			t.Fatalf("unknown health service: got %v, want NotFound", err)
		}
		return
	}
	if err != nil || response.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("health %s = %v, %v", address, response, err)
	}
}

func stop(t *testing.T, p *process, service string) {
	t.Helper()
	// Windows cannot send os.Interrupt to an arbitrary child. Its forced cleanup
	// verifies isolation; normal signal exit/logs are verified on both Linux CIs.
	var err error
	if runtime.GOOS == "windows" {
		err = p.cmd.Process.Kill()
	} else {
		err = p.cmd.Process.Signal(os.Interrupt)
	}
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-p.done:
	case <-time.After(10 * time.Second):
		t.Fatal("service did not stop")
	}
	records := p.logs.check(t, service)
	if runtime.GOOS != "windows" {
		if p.waitErr != nil || len(records) != 2 || records[1]["msg"] != "service stopped" || records[1]["stop_mode"] != "immediate" {
			t.Fatalf("unexpected normal exit: %v; logs: %v", p.waitErr, records)
		}
		logJSON(t, records[1])
	}
}

func TestServiceProcesses(t *testing.T) {
	binaries := os.Getenv("LAB_BIN_DIR")
	if binaries == "" {
		binaries = t.TempDir()
		build := exec.CommandContext(t.Context(), "go", "build", "-o", binaries+string(os.PathSeparator), "./cmd/...")
		build.Dir = "../.."
		if output, err := build.CombinedOutput(); err != nil {
			t.Fatalf("build service binaries: %v\n%s", err, output)
		}
	}
	services := []string{"orders", "inventory", "payments"}
	processes := make([]*process, len(services))
	for i, name := range services {
		binary := filepath.Join(binaries, name)
		if runtime.GOOS == "windows" {
			binary += ".exe"
		}
		processes[i] = start(t, binary)
		checkHealth(t, processes[i].address, "")
		checkHealth(t, processes[i].address, "unknown-business-api")
		records := processes[i].logs.check(t, name)
		for _, key := range []string{"version", "revision", "vcs_modified", "go_version"} {
			if records[0][key] == nil || records[0][key] == "" {
				t.Errorf("%s startup lacks %s", name, key)
			}
		}
		t.Logf("%s health SERVING", name)
		logJSON(t, records[0])
		for key, value := range map[string]string{"FL_GRPC_ADDR": "secret://user:password@host:80", "FL_LOG_LEVEL": "super-secret-value", "FL_SHUTDOWN_TIMEOUT": "0s"} {
			t.Run(name+"/reject/"+key, func(t *testing.T) {
				values := map[string]string{"FL_GRPC_ADDR": processes[i].address}
				values[key] = value
				assertFailure(t, binary, name, values, key, value)
			})
		}
		t.Run(name+"/missing-address", func(t *testing.T) {
			assertFailure(t, binary, name, nil, "FL_GRPC_ADDR is required", "")
		})
		t.Run(name+"/empty-address", func(t *testing.T) {
			assertFailure(t, binary, name, map[string]string{"FL_GRPC_ADDR": ""}, "FL_GRPC_ADDR is required", "")
		})
		t.Run(name+"/occupied", func(t *testing.T) {
			assertFailure(t, binary, name, map[string]string{"FL_GRPC_ADDR": processes[i].address}, "listen for gRPC", "")
		})
	}
	stop(t, processes[0], "orders")
	for i := 1; i < len(services); i++ {
		checkHealth(t, processes[i].address, "")
		t.Logf("%s remains SERVING after orders stopped", services[i])
	}
	for i := 1; i < len(services); i++ {
		stop(t, processes[i], services[i])
	}
}

func assertFailure(t *testing.T, binary, service string, overrides map[string]string, reason, secret string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Env = childEnvironment(overrides)
	output, err := cmd.CombinedOutput()
	var exitError *exec.ExitError
	if ctx.Err() != nil || !errors.As(err, &exitError) || exitError.ExitCode() != 1 {
		t.Fatalf("expected prompt exit 1, got %v; output %s", err, output)
	}
	var record map[string]any
	if err := json.Unmarshal(output, &record); err != nil {
		t.Fatalf("failure was not one JSON record: %v; %s", err, output)
	}
	if record["service"] != service || record["level"] != "ERROR" || !strings.Contains(fmt.Sprint(record["error"]), reason) || (secret != "" && strings.Contains(string(output), secret)) {
		t.Fatalf("unexpected failure record: %v", record)
	}
	logJSON(t, record)
}

func logJSON(t *testing.T, record map[string]any) {
	t.Helper()
	data, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("observed JSON: %s", data)
}
