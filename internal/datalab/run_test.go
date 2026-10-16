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
)

func TestResetRejectsAmbiguousTargets(t *testing.T) {
	for _, test := range []struct{ project, confirm string }{
		{"", ""}, {"faultline-lab", "faultline-lab"}, {"foreign", "foreign"},
		{"faultline-data-local", ""}, {"faultline-data-local", "faultline-data-other"},
		{"faultline-data-../other", "faultline-data-../other"},
	} {
		if err := Reset(context.Background(), test.project, test.confirm, io.Discard); err == nil {
			t.Fatalf("accepted unsafe reset target %q", test.project)
		}
	}
}

func TestRunCleanupAndReports(t *testing.T) {
	for _, scenario := range []string{"success", "read-only-directory", "partial-build", "cancel-build", "cleanup-failure"} {
		t.Run(scenario, func(t *testing.T) {
			root := t.TempDir()
			config, e := os.ReadFile("../../configs/postgres/local.json")
			if e != nil {
				t.Fatal(e)
			}
			if e = os.MkdirAll(filepath.Join(root, "configs/postgres"), 0o700); e != nil {
				t.Fatal(e)
			}
			if e = os.WriteFile(filepath.Join(root, "configs/postgres/local.json"), config, 0o600); e != nil {
				t.Fatal(e)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			cleaned := false
			imageExists := false
			execute := func(commandCtx context.Context, program string, args []string) (command.Output, error) {
				if program != "docker" {
					t.Fatal("unexpected executable")
				}
				if slices.Contains(args, "build") {
					imageExists = true
				}
				if args[0] == "image" && imageExists {
					return command.Output{Stdout: "owned-dbtool-image"}, nil
				}
				if slices.Contains(args, "build") && (scenario == "partial-build" || scenario == "cancel-build") {
					if scenario == "cancel-build" {
						cancel()
					}
					return command.Output{ExitCode: 1}, errors.New("controlled failure in " + root)
				}
				if slices.Contains(args, "down") {
					if scenario == "read-only-directory" {
						paths, e := filepath.Glob(filepath.Join(root, "artifacts/postgres/*/secrets"))
						if e != nil || len(paths) != 1 {
							t.Fatal("secret directory missing", e)
						}
						if e = os.Chmod(paths[0], 0o500); e != nil {
							t.Fatal(e)
						}
					}
					if commandCtx.Err() != nil {
						t.Fatal("cleanup reused canceled context")
					}
					if !slices.Contains(args, "--volumes") || !slices.Contains(args, "--rmi") || args[len(args)-1] != "local" {
						t.Fatal("wrong cleanup scope")
					}
					cleaned = true
					if scenario == "cleanup-failure" {
						return command.Output{ExitCode: 1}, errors.New("cleanup unavailable")
					}
					if slices.Contains(args, "--profile") && slices.Contains(args, "tools") {
						imageExists = false
					}
				}
				return command.Output{}, nil
			}
			err := run(ctx, io.Discard, root, execute)
			wantSuccess := scenario == "success" || scenario == "read-only-directory"
			if !cleaned || wantSuccess != (err == nil) {
				t.Fatalf("cleanup=%t err=%v", cleaned, err)
			}
			paths, e := filepath.Glob(filepath.Join(root, "artifacts/postgres/*/report.json"))
			if e != nil || len(paths) != 1 {
				t.Fatal("report missing", e)
			}
			data, e := os.ReadFile(paths[0])
			if e != nil {
				t.Fatal(e)
			}
			var saved report
			if e = json.Unmarshal(data, &saved); e != nil {
				t.Fatal(e)
			}
			if strings.Contains(string(data), root) {
				t.Fatal("report exposes host path")
			}
			if _, e = os.Stat(filepath.Join(filepath.Dir(paths[0]), "secrets")); !errors.Is(e, os.ErrNotExist) {
				t.Fatal("temporary credentials remain")
			}
			if wantSuccess && saved.Result != "PASS" || !wantSuccess && saved.Result != "FAIL" {
				t.Fatal("wrong verdict")
			}
			if (scenario == "cleanup-failure") != (saved.Cleanup == "FAIL") {
				t.Fatal("wrong cleanup verdict")
			}
		})
	}
}
