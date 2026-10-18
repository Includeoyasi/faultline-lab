// Package smoke owns a single isolated Docker Compose smoke run and its cleanup.
package smoke

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var services = []string{"orders", "inventory", "payments"}

// Options controls execution budgets, never service deployment defaults.
type Options struct {
	Root           string
	Artifacts      string
	Output         io.Writer
	CommandTimeout time.Duration
	BuildTimeout   time.Duration
	CleanupTimeout time.Duration
}

type runner struct {
	options           Options
	directory         string
	report            Report
	execute           executeFunc
	resourcesMayExist bool
	redact            *strings.Replacer
}

// Run invokes Docker directly on the host, without a shell or Docker socket mount.
func Run(ctx context.Context, options Options) (Report, error) {
	root, err := filepath.Abs(options.Root)
	if err != nil {
		return Report{}, errors.New("cannot resolve repository directory")
	}
	options.Root = root
	executor := processExecutor{directory: root, env: smokeEnvironment()}
	return run(ctx, options, executor.execute)
}

func smokeEnvironment() []string {
	env := make([]string, 0, len(os.Environ())+2)
	for _, value := range os.Environ() {
		key, _, _ := strings.Cut(value, "=")
		if strings.EqualFold(key, "FL_RUNTIME_ENV_FILE") || strings.EqualFold(key, "FL_STOP_GRACE_PERIOD") {
			continue
		}
		env = append(env, value)
	}
	return append(env, "FL_RUNTIME_ENV_FILE=configs/env/container.env.example", "FL_STOP_GRACE_PERIOD=15s")
}

func run(ctx context.Context, options Options, execute executeFunc) (Report, error) {
	r, err := newRunner(options, execute)
	if err != nil {
		return Report{}, err
	}
	if err = r.checkpoint(); err != nil {
		return r.report, errors.New("cannot save initial smoke report")
	}
	err = r.exercise(ctx)
	if err != nil {
		r.report.Result = "FAIL"
		if r.report.Phase == "prerequisites" && ctx.Err() == nil {
			r.report.Result = "NOT_RUN"
		}
		r.report.FailedPhase, r.report.Error = r.report.Phase, r.redact.Replace(err.Error())
	} else {
		r.report.Result = "PASS"
	}
	// The exercise has returned and joined its current child. Cleanup gets a new
	// budget: cancellation of the run must not cancel resource removal as well.
	cleanupCtx, cancel := context.WithTimeout(context.Background(), r.options.CleanupTimeout)
	cleanupErr := r.cleanup(cleanupCtx)
	cancel()
	if cleanupErr != nil {
		r.report.Result, r.report.Cleanup = "FAIL", "FAIL"
		r.report.CleanupError = r.redact.Replace(cleanupErr.Error())
	}
	if ctx.Err() != nil {
		r.report.Cancelled, r.report.Result = true, "FAIL"
		if r.report.Error == "" {
			r.report.Error = ctx.Err().Error()
		}
	}
	finished := time.Now().UTC()
	r.report.FinishedAt = &finished
	saveErr := r.checkpoint()
	if saveErr != nil {
		r.report.Result = "FAIL"
	}
	_, _ = fmt.Fprintf(r.options.Output, "%s: artifacts/smoke/%s/report.json\n", r.report.Result, r.report.RunID)
	return r.report, errors.Join(err, cleanupErr, saveErr)
}

func newRunner(options Options, execute executeFunc) (*runner, error) {
	if options.Output == nil {
		options.Output = io.Discard
	}
	if options.CommandTimeout <= 0 {
		options.CommandTimeout = time.Minute
	}
	if options.BuildTimeout <= 0 {
		options.BuildTimeout = 10 * time.Minute
	}
	if options.CleanupTimeout <= 0 {
		options.CleanupTimeout = time.Minute
	}
	if options.Artifacts == "" {
		options.Artifacts = filepath.Join(options.Root, "artifacts", "smoke")
	}
	var token [16]byte
	if _, err := rand.Read(token[:]); err != nil {
		return nil, errors.New("cannot create smoke run ID")
	}
	id := hex.EncodeToString(token[:])
	id = id[:8] + "-" + id[8:12] + "-" + id[12:16] + "-" + id[16:20] + "-" + id[20:]
	project := "faultline-smoke-" + id[:8]
	directory := filepath.Join(options.Artifacts, id)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return nil, errors.New("cannot create smoke artifact directory")
	}
	pairs := []string{options.Root, "[workspace]", filepath.ToSlash(options.Root), "[workspace]"}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		pairs = append(pairs, home, "[home]", filepath.ToSlash(home), "[home]")
	}
	return &runner{options: options, directory: directory, execute: execute, redact: strings.NewReplacer(pairs...), report: Report{
		SchemaVersion: 2, RunID: id, Project: project, StartedAt: time.Now().UTC(), Result: "FAIL", Phase: "prerequisites", Cleanup: "PENDING",
		HostPlatform: runtime.GOOS, GoVersion: runtime.Version(), Services: map[string]serviceReport{}, Commands: []commandReport{},
		Recovery:    []string{"docker", "compose", "-f", "compose.yaml", "-p", project, "down", "--remove-orphans", "--rmi", "local"},
		Limitations: []string{"No business operations or capacity measurement", "docker_architecture describes the backend, not individual image platforms (F03)", "SIGKILL or host failure requires recovery_command; signal cleanup cannot be guaranteed in those cases"},
	}}, nil
}

func (r *runner) command(ctx context.Context, program string, args []string, timeout time.Duration, sensitive bool) (commandOutput, error) {
	if err := ctx.Err(); err != nil {
		return commandOutput{}, err
	}
	commandCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	observation := commandReport{Program: program, Args: args, StartedAt: time.Now().UTC()}
	result, err := r.execute(commandCtx, program, args)
	observation.FinishedAt, observation.ExitCode = time.Now().UTC(), result.exitCode
	if err != nil {
		observation.Error = r.redact.Replace(err.Error())
		if commandCtx.Err() != nil {
			observation.Error = r.redact.Replace(errors.Join(commandCtx.Err(), err).Error())
		}
		if !sensitive {
			name := fmt.Sprintf("command-%d.log", len(r.report.Commands)+1)
			diagnostic := r.redact.Replace(result.stdout + result.stderr)
			if writeErr := os.WriteFile(filepath.Join(r.directory, name), []byte(diagnostic), 0o600); writeErr != nil {
				observation.Error += "; cannot save command diagnostics"
			}
		}
		err = fmt.Errorf("%s %s failed (exit %d; %s)", program, strings.Join(args, " "), result.exitCode, observation.Error)
	}
	r.report.Commands = append(r.report.Commands, observation)
	return result, errors.Join(err, r.checkpoint())
}

func (r *runner) docker(ctx context.Context, args ...string) (commandOutput, error) {
	return r.command(ctx, "docker", args, r.options.CommandTimeout, false)
}

func (r *runner) composeArgs(args ...string) []string {
	return append([]string{"compose", "-f", "compose.yaml", "-p", r.report.Project}, args...)
}

func (r *runner) compose(ctx context.Context, args ...string) (commandOutput, error) {
	return r.docker(ctx, r.composeArgs(args...)...)
}

func (r *runner) exercise(ctx context.Context) error {
	steps := []struct {
		name string
		work func(context.Context) error
	}{
		{"prerequisites", r.prerequisites},
		{"configuration", r.configuration},
		{"build", func(ctx context.Context) error {
			r.resourcesMayExist = true // A partial build can leave run-specific images.
			_, err := r.command(ctx, "docker", r.composeArgs(append([]string{"build"}, services...)...), r.options.BuildTimeout, false)
			return err
		}},
		{"start", func(ctx context.Context) error {
			r.resourcesMayExist = true
			_, err := r.compose(ctx, append([]string{"up", "-d", "--wait", "--wait-timeout", "60"}, services...)...)
			return err
		}},
		{"check", r.checkAll},
		{"isolation", r.isolate},
		{"stop", r.stopAll},
	}
	for _, step := range steps {
		if err := ctx.Err(); err != nil {
			return err
		}
		r.report.Phase = step.name
		if err := r.checkpoint(); err != nil {
			return errors.New("cannot checkpoint smoke phase")
		}
		_, _ = fmt.Fprintf(r.options.Output, "Runtime smoke %s: %s\n", r.report.Project, step.name)
		if err := step.work(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *runner) prerequisites(ctx context.Context) error {
	version, err := r.docker(ctx, "version", "--format", "{{.Server.Version}}")
	if err != nil {
		return err
	}
	r.report.DockerVersion = strings.TrimSpace(version.stdout)
	platform, err := r.docker(ctx, "info", "--format", "{{.OSType}}/{{.Architecture}}")
	if err != nil {
		return err
	}
	r.report.DockerArch = strings.TrimSpace(platform.stdout)
	if !strings.HasPrefix(r.report.DockerArch, "linux/") {
		return errors.New("docker must run Linux containers")
	}
	compose, err := r.docker(ctx, "compose", "version", "--short")
	if err != nil {
		return err
	}
	r.report.Compose = strings.TrimSpace(compose.stdout)
	revision, err := r.command(ctx, "git", []string{"rev-parse", "HEAD"}, r.options.CommandTimeout, false)
	if err != nil {
		return err
	}
	r.report.Revision = strings.TrimSpace(revision.stdout)
	dirty, err := r.command(ctx, "git", []string{"status", "--porcelain"}, r.options.CommandTimeout, true)
	r.report.Dirty = strings.TrimSpace(dirty.stdout) != ""
	return err
}

func (r *runner) configuration(ctx context.Context) error {
	resolved, err := r.command(ctx, "docker", r.composeArgs("config", "--format", "json"), r.options.CommandTimeout, true)
	if err != nil {
		return err
	}
	_, err = validateConfig([]byte(resolved.stdout), false)
	return err
}

func (r *runner) cleanup(ctx context.Context) error {
	if !r.resourcesMayExist {
		r.report.Cleanup = "NOT_NEEDED"
		return nil
	}
	for _, service := range services {
		logCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		_, err := r.collectLogs(logCtx, service)
		cancel()
		if err != nil {
			observation := r.report.Services[service]
			observation.DiagnosticError = r.redact.Replace(err.Error())
			r.report.Services[service] = observation
		}
	}
	// Recheck engine state after the joined up command and scoped down. A second
	// observation catches resources still being finalized by the Docker engine.
	for {
		if _, err := r.compose(ctx, "down", "--remove-orphans", "--rmi", "local"); err != nil {
			return err
		}
		for attempt := 0; attempt < 2; attempt++ {
			for _, args := range [][]string{
				{"ps", "--all", "--quiet", "--filter", "label=com.docker.compose.project=" + r.report.Project},
				{"network", "ls", "--quiet", "--filter", "label=com.docker.compose.project=" + r.report.Project},
				{"image", "ls", "--quiet", "--filter", "reference=" + r.report.Project + "-*:latest"},
			} {
				remaining, err := r.docker(ctx, args...)
				if err != nil {
					return err
				}
				if strings.TrimSpace(remaining.stdout) != "" {
					goto retry
				}
			}
			if attempt == 0 {
				timer := time.NewTimer(200 * time.Millisecond)
				select {
				case <-ctx.Done():
					timer.Stop()
					return ctx.Err()
				case <-timer.C:
				}
			}
		}
		r.report.Cleanup = "PASS"
		return nil
	retry:
		if err := ctx.Err(); err != nil {
			return err
		}
		timer := time.NewTimer(200 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}
