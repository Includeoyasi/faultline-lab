package datalab

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/Includeoyasi/faultline-lab/internal/command"
	"github.com/Includeoyasi/faultline-lab/internal/infrastructure"
)

func TestCombinedInfrastructureLifecycle(t *testing.T) {
	for _, scenario := range []string{"success", "start-failure", "before-failure", "invalid-checkpoint", "after-failure", "cancel-before", "cleanup-failure", "metadata-cancel", "metadata-missing"} {
		t.Run(scenario, func(t *testing.T) {
			root := t.TempDir()
			for _, path := range []string{"configs/postgres/local.json", "configs/infrastructure/local.json"} {
				data, err := os.ReadFile(filepath.Join("../..", path))
				if err != nil {
					t.Fatal(err)
				}
				destination := filepath.Join(root, path)
				if err = os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
					t.Fatal(err)
				}
				if err = os.WriteFile(destination, data, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			cleaned, logged, afterCalled := false, false, false
			var password, passwordHash string
			checkpoint := `{"version":1,"topic":"faultline-smoke","partition":0,"offset":7,"key":"checkpoint","value":"expected"}`
			execute := func(commandCtx context.Context, program string, args []string) (command.Output, error) {
				if program != "docker" {
					t.Fatal("unexpected executable")
				}
				if scenario == "metadata-cancel" && args[0] == "info" {
					cancel()
					return command.Output{}, commandCtx.Err()
				}
				if scenario == "metadata-missing" && args[0] == "image" && args[1] == "inspect" {
					return command.Output{Stderr: "metadata-private"}, errors.New("metadata-private")
				}
				if data, ok := metadataFixture(args); ok {
					return data, nil
				}
				if args[0] == "compose" && !slices.Contains(args, "compose.infrastructure.yaml") {
					t.Fatal("combined command lost its overlay")
				}
				if slices.Contains(args, "logs") {
					logged = true
					return command.Output{Stdout: "diagnostic " + password + " " + passwordHash + " " + root}, nil
				}
				if slices.Contains(args, "build") {
					paths, err := filepath.Glob(filepath.Join(root, "artifacts/infrastructure/*/infrastructure-secrets"))
					if err != nil || len(paths) != 1 {
						t.Fatal("missing isolated infrastructure credentials")
					}
					cfg, err := infrastructure.Load(filepath.Join(root, "configs/infrastructure/local.json"), paths[0])
					if err != nil {
						t.Fatal(err)
					}
					password, err = infrastructure.Password(cfg)
					if err != nil {
						t.Fatal(err)
					}
					passwordHash = infrastructure.PasswordHash(password)
				}
				if slices.Contains(args, "down") {
					if commandCtx.Err() != nil || !slices.Contains(args, "--volumes") || !slices.Contains(args, "--profile") {
						t.Fatal("cleanup has incorrect scope or canceled context")
					}
					cleaned = true
					if scenario == "cleanup-failure" {
						return command.Output{ExitCode: 1}, errors.New("controlled cleanup failure")
					}
				}
				if slices.Contains(args, "up") {
					for _, name := range []string{"postgres", "kafka", "redis"} {
						if !slices.Contains(args, name) {
							t.Fatal("combined start omitted a dependency")
						}
					}
					if scenario == "start-failure" {
						return command.Output{ExitCode: 1}, errors.New("controlled startup failure")
					}
				}
				if slices.Contains(args, "before-restart") {
					if scenario == "invalid-checkpoint" {
						return command.Output{Stdout: `{"version":0}`}, nil
					}
					if scenario == "before-failure" || scenario == "cancel-before" {
						if scenario == "cancel-before" {
							cancel()
						}
						return command.Output{ExitCode: 1}, errors.New("controlled check failure")
					}
					return command.Output{Stdout: checkpoint}, nil
				}
				if slices.Contains(args, "after-restart") {
					afterCalled = true
					index := slices.Index(args, "-checkpoint")
					if index < 0 || args[index+1] != checkpoint {
						t.Fatal("after-restart did not reuse the saved checkpoint")
					}
					if scenario == "after-failure" {
						return command.Output{ExitCode: 1}, errors.New("controlled re-read failure")
					}
				}
				return command.Output{}, nil
			}
			err := runMode(ctx, io.Discard, root, execute, true)
			wantSuccess := scenario == "success" || scenario == "metadata-missing"
			if !cleaned || !logged || (err == nil) != wantSuccess {
				t.Fatal("incorrect combined lifecycle result")
			}
			if (scenario == "start-failure" || scenario == "before-failure" || scenario == "cancel-before" || scenario == "invalid-checkpoint" || scenario == "metadata-cancel") && afterCalled {
				t.Fatal("continued after an earlier failure")
			}
			paths, err := filepath.Glob(filepath.Join(root, "artifacts/infrastructure/*/report.json"))
			if err != nil || len(paths) != 1 {
				t.Fatal("combined report missing")
			}
			data, err := os.ReadFile(paths[0])
			if err != nil || strings.Contains(string(data), root) {
				t.Fatal("combined report unavailable or exposes host path")
			}
			if password == "" || strings.Contains(string(data), password) || strings.Contains(string(data), passwordHash) || !strings.Contains(string(data), "[secret hash]") {
				t.Fatal("report did not redact actual credential output")
			}
			var saved report
			if json.Unmarshal(data, &saved) != nil || (saved.Result == "PASS") != wantSuccess {
				t.Fatal("incorrect saved verdict")
			}
			if scenario == "success" {
				for _, service := range []string{"postgres", "kafka", "redis"} {
					if saved.Environment.Services[service].Image.Status != "KNOWN" || saved.Environment.Services[service].Usage.Status != "MEASURED" {
						t.Fatal("combined metadata missing", service)
					}
				}
			}
			if strings.Contains(string(data), "metadata-secret") || strings.Contains(string(data), strings.Join([]string{"", "Users", "private"}, "/")) {
				t.Fatal("raw Compose data leaked")
			}
			if strings.Contains(string(data), "metadata-private") || scenario == "metadata-missing" && saved.Environment.Services["postgres"].Image.Status != "UNKNOWN" {
				t.Fatal("missing metadata not safely represented")
			}
			for _, name := range []string{"secrets", "infrastructure-secrets"} {
				if _, err = os.Stat(filepath.Join(filepath.Dir(paths[0]), name)); !errors.Is(err, os.ErrNotExist) {
					t.Fatal("temporary credentials remain")
				}
			}
		})
	}
}
