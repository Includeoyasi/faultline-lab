package smoke

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var containerIDPattern = regexp.MustCompile(`^[a-f0-9]{12,64}$`)

type container struct {
	Image  string
	Config struct{ User string }
	State  struct {
		Status    string
		ExitCode  int
		OOMKilled bool
		Health    struct{ Status string }
	}
	HostConfig struct {
		ReadonlyRootfs bool
		Privileged     bool
		CapDrop        []string
		SecurityOpt    []string
		PortBindings   map[string]json.RawMessage
	}
	Mounts          []json.RawMessage
	NetworkSettings struct {
		Networks map[string]struct{ NetworkID string }
	}
}

func (r *runner) inspect(ctx context.Context, service string) (string, container, error) {
	ids, err := r.compose(ctx, "ps", "--all", "-q", service)
	if err != nil {
		return "", container{}, err
	}
	id := strings.TrimSpace(ids.stdout)
	if !containerIDPattern.MatchString(id) {
		return "", container{}, fmt.Errorf("missing unique container for %s", service)
	}
	result, err := r.docker(ctx, "inspect", id)
	if err != nil {
		return "", container{}, err
	}
	var data []container
	if err = json.Unmarshal([]byte(result.stdout), &data); err != nil || len(data) != 1 {
		return "", container{}, fmt.Errorf("invalid container inspect for %s", service)
	}
	return id, data[0], nil
}

func (r *runner) checkAll(ctx context.Context) error {
	for _, service := range services {
		if err := r.checkRunning(ctx, service); err != nil {
			return err
		}
	}
	return nil
}

func (r *runner) checkRunning(ctx context.Context, service string) error {
	if _, err := r.compose(ctx, "exec", "-T", service, "/app/service", "healthcheck"); err != nil {
		return err
	}
	id, data, err := r.inspect(ctx, service)
	if err != nil {
		return err
	}
	host := data.HostConfig
	observation := r.report.Services[service]
	observation.ID, observation.ImageID, observation.Status = id, data.Image, data.State.Status
	observation.OOMKilled = &data.State.OOMKilled
	r.report.Services[service] = observation
	valid := data.State.Status == "running" && data.State.Health.Status == "healthy" &&
		data.Config.User == "65532:65532" && host.ReadonlyRootfs && !host.Privileged &&
		slices.Contains(host.CapDrop, "ALL") && slices.ContainsFunc(host.SecurityOpt, func(v string) bool {
		return v == "no-new-privileges" || v == "no-new-privileges:true" || v == "no-new-privileges=true"
	}) &&
		len(host.PortBindings) == 0 && len(data.Mounts) == 0 && len(data.NetworkSettings.Networks) > 0
	if !valid {
		return fmt.Errorf("%s violates health or runtime security requirements", service)
	}
	for _, network := range data.NetworkSettings.Networks {
		result, err := r.docker(ctx, "network", "inspect", network.NetworkID)
		if err != nil {
			return err
		}
		var networks []struct{ Internal bool }
		if err = json.Unmarshal([]byte(result.stdout), &networks); err != nil || len(networks) != 1 || !networks[0].Internal {
			return fmt.Errorf("%s must use an internal network", service)
		}
	}
	observation.ID, observation.ImageID, observation.Healthy = id, data.Image, true
	observation.Security = map[string]any{"user": data.Config.User, "read_only": true, "capabilities_dropped": true, "no_new_privileges": true, "host_ports": []string{}, "mounts": []string{}, "internal_network": true}
	r.report.Services[service] = observation
	return nil
}

func (r *runner) collectLogs(ctx context.Context, service string) (string, error) {
	result, err := r.compose(ctx, "logs", "--no-color", "--no-log-prefix", service)
	if err != nil {
		return "", err
	}
	logs := r.redact.Replace(result.stdout + result.stderr)
	if err = os.WriteFile(filepath.Join(r.directory, service+".jsonl"), []byte(logs), 0o600); err != nil {
		return "", errors.New("cannot save service logs")
	}
	return logs, nil
}

func (r *runner) checkStopped(ctx context.Context, service string) error {
	id, data, err := r.inspect(ctx, service)
	if err != nil {
		return err
	}
	observation := r.report.Services[service]
	observation.ID, observation.ImageID, observation.Status = id, data.Image, data.State.Status
	observation.ExitCode, observation.OOMKilled = &data.State.ExitCode, &data.State.OOMKilled
	r.report.Services[service] = observation
	if data.State.Status != "exited" || data.State.ExitCode != 0 || data.State.OOMKilled {
		return fmt.Errorf("%s failed normal SIGTERM shutdown", service)
	}
	logs, err := r.collectLogs(ctx, service)
	if err != nil {
		return err
	}
	step := 0
	for _, line := range strings.Split(logs, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var record struct {
			Service, Time, Level, Msg, Health string
			StopMode                          string `json:"stop_mode"`
		}
		if err = json.Unmarshal([]byte(line), &record); err != nil || record.Service != service || record.Time == "" || record.Level == "" || record.Msg == "" {
			return fmt.Errorf("%s emitted an invalid JSON record", service)
		}
		switch {
		case step == 0 && record.Msg == "service started":
			step++
		case step == 1 && record.Msg == "service stopping" && record.Health == "NOT_SERVING":
			step++
		case step == 2 && record.Msg == "service stopped" && record.StopMode == "graceful":
			step++
		}
	}
	if step != 3 {
		return fmt.Errorf("%s lacks ordered start/readiness/stop records", service)
	}
	observation.Stopped, observation.ExitCode, observation.StopMode = true, &data.State.ExitCode, "graceful"
	r.report.Services[service] = observation
	return nil
}

func (r *runner) isolate(ctx context.Context) error {
	if _, err := r.compose(ctx, "stop", "orders"); err != nil {
		return err
	}
	if err := r.checkStopped(ctx, "orders"); err != nil {
		return err
	}
	for _, service := range []string{"inventory", "payments"} {
		if err := r.checkRunning(ctx, service); err != nil {
			return err
		}
	}
	r.report.Isolation = "PASS"
	return nil
}

func (r *runner) stopAll(ctx context.Context) error {
	if _, err := r.compose(ctx, "stop", "inventory", "payments"); err != nil {
		return err
	}
	for _, service := range []string{"inventory", "payments"} {
		if err := r.checkStopped(ctx, service); err != nil {
			return err
		}
	}
	return nil
}
