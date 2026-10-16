package smoke

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Report version 2 preserves result/cleanup/services/commands and exit meanings.
// It replaces Node metadata with Go metadata and adds cancellation/recovery data.
type Report struct {
	SchemaVersion int                      `json:"schema_version"`
	RunID         string                   `json:"run_id"`
	Project       string                   `json:"project"`
	StartedAt     time.Time                `json:"started_at"`
	FinishedAt    *time.Time               `json:"finished_at,omitempty"`
	Result        string                   `json:"result"`
	Phase         string                   `json:"phase"`
	FailedPhase   string                   `json:"failed_phase,omitempty"`
	Error         string                   `json:"error,omitempty"`
	Cancelled     bool                     `json:"cancelled"`
	Cleanup       string                   `json:"cleanup"`
	CleanupError  string                   `json:"cleanup_error,omitempty"`
	HostPlatform  string                   `json:"host_platform"`
	GoVersion     string                   `json:"go_version"`
	DockerVersion string                   `json:"docker_version,omitempty"`
	DockerArch    string                   `json:"docker_architecture,omitempty"`
	Compose       string                   `json:"compose_version,omitempty"`
	Revision      string                   `json:"source_revision,omitempty"`
	Dirty         bool                     `json:"source_dirty"`
	Isolation     string                   `json:"isolation_after_orders_stop,omitempty"`
	Services      map[string]serviceReport `json:"services"`
	Commands      []commandReport          `json:"commands"`
	Recovery      []string                 `json:"recovery_command"`
	Limitations   []string                 `json:"limitations"`
}

type serviceReport struct {
	ID              string         `json:"id,omitempty"`
	ImageID         string         `json:"image_id,omitempty"`
	Healthy         bool           `json:"healthy"`
	Stopped         bool           `json:"stopped"`
	ExitCode        *int           `json:"exit_code,omitempty"`
	OOMKilled       *bool          `json:"oom_killed,omitempty"`
	Status          string         `json:"status,omitempty"`
	StopMode        string         `json:"stop_mode,omitempty"`
	Security        map[string]any `json:"security,omitempty"`
	DiagnosticError string         `json:"diagnostic_error,omitempty"`
}

type commandReport struct {
	Program    string    `json:"program"`
	Args       []string  `json:"args"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	ExitCode   int       `json:"exit_code"`
	Error      string    `json:"error,omitempty"`
}

// ExitCode keeps the original smoke's public exit contract.
func (r Report) ExitCode() int {
	switch r.Result {
	case "PASS":
		return 0
	case "NOT_RUN":
		return 2
	default:
		return 1
	}
}

func (r *runner) checkpoint() error {
	data, err := json.MarshalIndent(r.report, "", "  ")
	if err != nil {
		return err
	}
	temporary := filepath.Join(r.directory, "report.json.tmp")
	if err = os.WriteFile(temporary, append(data, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(temporary, filepath.Join(r.directory, "report.json"))
}
