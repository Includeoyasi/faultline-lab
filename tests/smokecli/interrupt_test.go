//go:build !windows

// Package smokecli runs opt-in host integration tests of the compiled CLI.
package smokecli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/smoke"
)

func TestSIGINTCleansRunningProject(t *testing.T) {
	binary, root := os.Getenv("LAB_SMOKE_BINARY"), os.Getenv("LAB_SMOKE_ROOT")
	if binary == "" || root == "" {
		t.Skip("host integration requires LAB_SMOKE_BINARY and LAB_SMOKE_ROOT; CI runs the compiled test on the host")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, "-root", root)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}
	projects, signals, scanned := make(chan string, 1), make(chan error, 1), make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasSuffix(line, ": prerequisites") {
				fields := strings.Fields(line)
				if len(fields) >= 3 {
					projects <- strings.TrimSuffix(fields[2], ":")
				}
			}
			if strings.HasSuffix(line, ": check") {
				signals <- cmd.Process.Signal(os.Interrupt)
			}
		}
		scanned <- scanner.Err()
	}()
	scanErr := <-scanned
	waitErr := cmd.Wait()
	if scanErr != nil {
		t.Logf("CLI stdout ended with: %v", scanErr)
	}
	var project string
	select {
	case project = <-projects:
	default:
		t.Fatalf("CLI never identified its project: %v; %s", waitErr, stderr.String())
	}
	// Even a failed regression must leave no resources owned by this test.
	t.Cleanup(func() {
		cleanupCtx, stop := context.WithTimeout(context.Background(), time.Minute)
		defer stop()
		cleanup := exec.CommandContext(cleanupCtx, "docker", "compose", "-f", filepath.Join(root, "compose.yaml"), "-p", project, "down", "--remove-orphans")
		if output, err := cleanup.CombinedOutput(); err != nil {
			t.Errorf("test recovery failed: %v; %s", err, output)
		}
	})
	select {
	case err := <-signals:
		if err != nil {
			t.Fatal(err)
		}
	default:
		t.Fatalf("CLI did not reach interruption phase: %v; %s", waitErr, stderr.String())
	}
	var exit *exec.ExitError
	if !errors.As(waitErr, &exit) || exit.ExitCode() != 1 {
		t.Fatalf("unexpected interrupt exit: %v; %s", waitErr, stderr.String())
	}
	paths, err := filepath.Glob(filepath.Join(root, "artifacts", "smoke", "*", "report.json"))
	if err != nil {
		t.Fatal(err)
	}
	var observed *smoke.Report
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var report smoke.Report
		if err = json.Unmarshal(data, &report); err != nil {
			t.Fatal(err)
		}
		if report.Project == project {
			observed = &report
		}
	}
	if observed == nil || !observed.Cancelled || observed.Cleanup != "PASS" || observed.Result != "FAIL" || observed.FinishedAt == nil {
		t.Fatalf("missing or invalid cancellation report: %+v", observed)
	}
	for _, args := range [][]string{
		{"ps", "--all", "--quiet", "--filter", "label=com.docker.compose.project=" + project},
		{"network", "ls", "--quiet", "--filter", "label=com.docker.compose.project=" + project},
	} {
		output, err := exec.CommandContext(ctx, "docker", args...).Output()
		if err != nil || strings.TrimSpace(string(output)) != "" {
			t.Fatalf("owned resources remain: %s; %v", output, err)
		}
	}
}
