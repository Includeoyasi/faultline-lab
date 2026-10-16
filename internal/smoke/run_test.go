package smoke

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// dockerFixture models observable engine state and one independent control
// project. It tests ownership and failure paths without pretending to be Docker.
type dockerFixture struct {
	t                 *testing.T
	fail, cancelAt    string
	cancel            context.CancelFunc
	project           string
	owned, foreign    bool
	stopped           map[string]bool
	operations        []string
	unsafe, oom       bool
	cleanupBlocks     bool
	lateResourceReads int
}

func (f *dockerFixture) execute(ctx context.Context, program string, args []string) (commandOutput, error) {
	f.t.Helper()
	if program == "git" {
		if args[0] == "rev-parse" {
			return commandOutput{stdout: strings.Repeat("a", 40)}, nil
		}
		return commandOutput{}, nil
	}
	if program != "docker" {
		f.t.Fatalf("unexpected executable %s", program)
	}
	op := args[0]
	if op == "compose" && len(args) > 5 {
		project := args[4]
		if !strings.HasPrefix(project, "faultline-smoke-") || (f.project != "" && f.project != project) {
			f.t.Fatalf("unexpected project %s", project)
		}
		f.project = project
		args = args[5:]
		op = args[0]
	}
	f.operations = append(f.operations, op)
	if op == "up" {
		f.owned = true
	}
	if op == f.cancelAt {
		f.cancel()
		<-ctx.Done()
		return commandOutput{exitCode: -1}, ctx.Err()
	}
	if op == f.fail {
		return commandOutput{exitCode: 42, stderr: "controlled command failure"}, errors.New("fixture failure")
	}
	switch op {
	case "version":
		return commandOutput{stdout: "29.8.2"}, nil
	case "info":
		return commandOutput{stdout: "linux/arm64"}, nil
	case "compose": // compose version
		return commandOutput{stdout: "5.6.0"}, nil
	case "config":
		return commandOutput{stdout: string(resolvedConfig("10s", "15s"))}, nil
	case "build", "up", "exec":
		return commandOutput{}, nil
	case "stop":
		for _, name := range args[1:] {
			f.stopped[name] = true
		}
		return commandOutput{}, nil
	case "ps":
		if slices.Contains(args, "-q") {
			index := slices.Index(services, args[len(args)-1])
			return commandOutput{stdout: strings.Repeat(string(rune('a'+index)), 64)}, nil
		}
		if f.owned {
			return commandOutput{stdout: strings.Repeat("a", 64)}, nil
		}
		if f.lateResourceReads > 0 {
			f.lateResourceReads--
			f.owned = true
			return commandOutput{stdout: strings.Repeat("a", 64)}, nil
		}
		return commandOutput{}, nil
	case "inspect":
		index := int(args[1][0] - 'a')
		name := services[index]
		state, code, oom := "running", 0, false
		if f.stopped[name] {
			state = "exited"
			if f.oom {
				code, oom = 137, true
			}
		}
		user := "65532:65532"
		if f.unsafe {
			user = "root"
		}
		data := fmt.Sprintf(`[{"Image":"sha256:%s","Config":{"User":%q},"State":{"Status":%q,"ExitCode":%d,"OOMKilled":%t,"Health":{"Status":"healthy"}},"HostConfig":{"ReadonlyRootfs":true,"CapDrop":["ALL"],"SecurityOpt":["no-new-privileges:true"]},"Mounts":[],"NetworkSettings":{"Networks":{"runtime":{"NetworkID":"fixture-network"}}}}]`, strings.Repeat("d", 64), user, state, code, oom)
		return commandOutput{stdout: data}, nil
	case "network":
		if args[1] == "inspect" {
			return commandOutput{stdout: `[{"Internal":true}]`}, nil
		}
		return commandOutput{}, nil
	case "logs":
		name := args[len(args)-1]
		return commandOutput{stdout: fmt.Sprintf("{\"service\":%q,\"time\":\"2026-10-16T00:00:00Z\",\"level\":\"INFO\",\"msg\":\"service started\"}\n{\"service\":%q,\"time\":\"2026-10-16T00:00:01Z\",\"level\":\"INFO\",\"msg\":\"service stopping\",\"health\":\"NOT_SERVING\"}\n{\"service\":%q,\"time\":\"2026-10-16T00:00:02Z\",\"level\":\"INFO\",\"msg\":\"service stopped\",\"stop_mode\":\"graceful\"}\n", name, name, name)}, nil
	case "down":
		if f.cleanupBlocks {
			<-ctx.Done()
			return commandOutput{exitCode: -1}, ctx.Err()
		}
		f.owned = false
		return commandOutput{}, nil
	default:
		f.t.Fatalf("unexpected Docker operation %s %v", op, args)
		return commandOutput{}, nil
	}
}

func fixture(t *testing.T) (*dockerFixture, Options) {
	t.Helper()
	return &dockerFixture{t: t, foreign: true, stopped: map[string]bool{}}, Options{Root: t.TempDir(), CleanupTimeout: 2 * time.Second}
}

func verifySavedReport(t *testing.T, options Options, report Report) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(options.Root, "artifacts", "smoke", report.RunID, "report.json"))
	if err != nil {
		t.Fatal(err)
	}
	var saved Report
	if err = json.Unmarshal(data, &saved); err != nil || saved.Result != report.Result || saved.Cleanup != report.Cleanup || saved.Cancelled != report.Cancelled || saved.FinishedAt == nil {
		t.Fatalf("saved report mismatch: %s; %v", data, err)
	}
	if strings.Contains(string(data), options.Root) || len(saved.Recovery) == 0 {
		t.Fatal("report exposes workspace or lacks recovery command")
	}
}

func TestSuccessfulRunRetainsServiceChecks(t *testing.T) {
	f, options := fixture(t)
	f.lateResourceReads = 1
	report, err := run(context.Background(), options, f.execute)
	if err != nil || report.ExitCode() != 0 || report.Cleanup != "PASS" || f.owned || !f.foreign || report.Isolation != "PASS" {
		t.Fatalf("report=%+v err=%v", report, err)
	}
	for _, service := range services {
		observed := report.Services[service]
		if !observed.Healthy || !observed.Stopped || observed.ID == "" || observed.ImageID == "" || observed.ExitCode == nil || *observed.ExitCode != 0 {
			t.Fatalf("incomplete %s observation: %+v", service, observed)
		}
	}
	if count := strings.Count(strings.Join(f.operations, ","), "down"); count != 2 {
		t.Fatalf("late engine resource was not cleaned again: %v", f.operations)
	}
	verifySavedReport(t, options, report)
}

func TestCancellationStopsStagesAndCleansOwnProject(t *testing.T) {
	for _, op := range []string{"build", "up", "exec", "stop"} {
		t.Run(op, func(t *testing.T) {
			f, options := fixture(t)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			f.cancelAt, f.cancel = op, cancel
			report, err := run(ctx, options, f.execute)
			if err == nil || !report.Cancelled || report.ExitCode() != 1 || f.owned || !f.foreign {
				t.Fatalf("cancel report=%+v err=%v", report, err)
			}
			if op == "build" && report.Cleanup != "NOT_NEEDED" || op != "build" && report.Cleanup != "PASS" {
				t.Fatalf("cleanup=%s", report.Cleanup)
			}
			if op == "up" && slices.Contains(f.operations, "exec") {
				t.Fatal("new stage started after cancellation")
			}
			verifySavedReport(t, options, report)
		})
	}
}

func TestCommandFailuresAndCleanupDeadline(t *testing.T) {
	for _, op := range []string{"version", "build", "exec", "stop", "down", "cleanup-timeout"} {
		t.Run(op, func(t *testing.T) {
			f, options := fixture(t)
			f.fail = op
			if op == "cleanup-timeout" {
				f.cleanupBlocks, options.CleanupTimeout = true, 30*time.Millisecond
			}
			report, err := run(context.Background(), options, f.execute)
			wantCode := 1
			if op == "version" {
				wantCode = 2
			}
			if err == nil || report.ExitCode() != wantCode || !f.foreign {
				t.Fatalf("report=%+v err=%v", report, err)
			}
			if op == "down" || op == "cleanup-timeout" {
				if report.Cleanup != "FAIL" || !f.owned {
					t.Fatal("cleanup failure was hidden")
				}
			} else if f.owned {
				t.Fatal("owned resources left behind")
			}
			verifySavedReport(t, options, report)
		})
	}
}

func TestFailedChecksRetainObservedFacts(t *testing.T) {
	for _, scenario := range []string{"unsafe", "oom"} {
		t.Run(scenario, func(t *testing.T) {
			f, options := fixture(t)
			f.unsafe, f.oom = scenario == "unsafe", scenario == "oom"
			report, err := run(context.Background(), options, f.execute)
			observed := report.Services["orders"]
			if err == nil || report.ExitCode() != 1 || observed.ID == "" || observed.ImageID == "" {
				t.Fatalf("missing failed inspection: %+v %v", observed, err)
			}
			if f.oom && (observed.ExitCode == nil || *observed.ExitCode != 137 || observed.OOMKilled == nil || !*observed.OOMKilled) {
				t.Fatalf("OOM observations lost: %+v", observed)
			}
			verifySavedReport(t, options, report)
		})
	}
}

func TestCommandDiagnosticsRetainRedactedCause(t *testing.T) {
	for _, sensitive := range []bool{false, true} {
		t.Run(fmt.Sprint(sensitive), func(t *testing.T) {
			_, options := fixture(t)
			r, err := newRunner(options, func(context.Context, string, []string) (commandOutput, error) {
				return commandOutput{exitCode: -1, stdout: "resolved-config-must-not-leak"}, fmt.Errorf("cannot create process in %s: executable missing", options.Root)
			})
			if err != nil {
				t.Fatal(err)
			}
			_, err = r.command(context.Background(), "docker", []string{"compose", "config"}, time.Second, sensitive)
			if err == nil || !strings.Contains(err.Error(), "executable missing") || strings.Contains(err.Error(), options.Root) {
				t.Fatalf("cause lost or not redacted: %v", err)
			}
			if !strings.Contains(r.report.Commands[0].Error, "[workspace]") {
				t.Fatal("report lost redacted process cause")
			}
			if sensitive {
				if _, err = os.Stat(filepath.Join(r.directory, "command-1.log")); !errors.Is(err, os.ErrNotExist) {
					t.Fatal("sensitive output persisted")
				}
			}
		})
	}
}
