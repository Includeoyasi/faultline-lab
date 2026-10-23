package smokecli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/runmeta"
)

// These opt-in checks read completed real reports. Unlike the functional CLI
// exit code, they require observable metadata in the known CI deployment.
func TestSuccessfulReportsHaveObservedMetadata(t *testing.T) {
	root := os.Getenv("LAB_SMOKE_ROOT")
	if root == "" {
		t.Skip("requires completed runtime, PostgreSQL and infrastructure reports")
	}
	for _, test := range []struct {
		kind     string
		version  int
		services []string
	}{
		{"smoke", 3, []string{"orders", "inventory", "payments"}},
		{"postgres", 1, []string{"postgres"}},
		{"infrastructure", 1, []string{"postgres", "kafka", "redis"}},
	} {
		t.Run(test.kind, func(t *testing.T) {
			r := latestSuccessfulReport(t, root, test.kind)
			if r.SchemaVersion != test.version || r.Cleanup != "PASS" || r.Environment.SchemaVersion != 1 {
				t.Fatal("report schema/cleanup missing")
			}
			env := r.Environment
			if env.Runner.Status != "KNOWN" || env.Backend.Status != "KNOWN" || env.BackendBudget.CPUs.Status != "KNOWN" || env.BackendBudget.MemoryBytes.Status != "KNOWN" {
				t.Fatal("runner/backend observations missing")
			}
			if env.ExecutionProfile.Status != "UNVERIFIED" || env.ComparisonEligibility.Status != "NOT_EVALUATED" {
				t.Fatal("functional report claimed performance eligibility")
			}
			for _, name := range test.services {
				s, exists := env.Services[name]
				if !exists || s.ContainerID == "" || s.ImageID == "" || s.Image.Status != "KNOWN" {
					t.Fatalf("%s image provenance missing: %+v", name, s.Image)
				}
				if s.Usage.Status != "MEASURED" || s.Usage.ObservedAt == nil || s.Usage.MemoryBytes == nil || s.Usage.CPUPercent == nil {
					t.Fatalf("%s resource sample missing", name)
				}
				for _, cap := range []runmeta.Quantity{s.Limits.MemoryBytes, s.Limits.NanoCPUs} {
					if cap.Value == nil || cap.Status != "KNOWN" && cap.Status != "NOT_CONFIGURED" {
						t.Fatalf("%s configured resource field missing", name)
					}
				}
			}
		})
	}
}

type savedReport struct {
	SchemaVersion int            `json:"schema_version"`
	Started       time.Time      `json:"started_at"`
	Result        string         `json:"result"`
	Cleanup       string         `json:"cleanup"`
	Environment   runmeta.Report `json:"environment"`
}

func latestSuccessfulReport(t *testing.T, root, kind string) savedReport {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(root, "artifacts", kind, "*", "report.json"))
	if err != nil {
		t.Fatal(err)
	}
	var latest savedReport
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		var r savedReport
		if err = json.Unmarshal(data, &r); err != nil {
			t.Fatal(err)
		}
		if r.Result == "PASS" && r.Started.After(latest.Started) {
			latest = r
		}
	}
	if latest.Result != "PASS" {
		t.Fatal("no completed successful report", kind)
	}
	return latest
}
