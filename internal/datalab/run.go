// Package datalab runs checks on its own disposable Compose project.
package datalab

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/command"
	"github.com/Includeoyasi/faultline-lab/internal/database"
	"github.com/Includeoyasi/faultline-lab/internal/infrastructure"
	"github.com/Includeoyasi/faultline-lab/internal/runmeta"
)

var projectName = regexp.MustCompile(`^faultline-data-[a-z0-9][a-z0-9-]{0,47}$`)

func composeModeArgs(project string, fullStack bool, args ...string) []string {
	// Profile activation by an earlier explicit run is not retained by down.
	// All starts/builds below retain explicit service targets.
	base := []string{"compose", "-f", "compose.yaml", "-f", "compose.postgres.yaml"}
	if fullStack {
		base = append(base, "-f", "compose.infrastructure.yaml")
	}
	return append(append(base, "-p", project, "--profile", "tools"), args...)
}

// Reset requires an exact explicit project confirmation; ordinary down preserves data.
func Reset(ctx context.Context, project, confirm string, out io.Writer) error {
	return reset(ctx, project, confirm, out, false)
}

// ResetInfrastructure includes Kafka's persistent volume in the explicit reset.
func ResetInfrastructure(ctx context.Context, project, confirm string, out io.Writer) error {
	return reset(ctx, project, confirm, out, true)
}

func reset(ctx context.Context, project, confirm string, out io.Writer, fullStack bool) error {
	if !projectName.MatchString(project) || project != confirm {
		return errors.New("reset requires a faultline-data-* project and identical -confirm value")
	}
	root, err := os.Getwd()
	if err != nil {
		return errors.New("cannot locate checkout")
	}
	executor := command.Executor{Directory: root, Env: os.Environ()}
	bounded, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	result, err := executor.Execute(bounded, "docker", composeModeArgs(project, fullStack, "down", "--volumes", "--remove-orphans", "--rmi", "local"))
	if err != nil {
		return errors.New("scoped reset failed; inspect Docker for the named project")
	}
	_, err = io.WriteString(out, result.Stdout)
	return err
}

type event struct {
	Stage    string `json:"stage"`
	Output   string `json:"output,omitempty"`
	Passed   bool   `json:"passed"`
	Error    string `json:"error,omitempty"`
	ExitCode int    `json:"exit_code"`
}
type report struct {
	SchemaVersion int            `json:"schema_version"`
	Project       string         `json:"project"`
	Started       time.Time      `json:"started_at"`
	Finished      time.Time      `json:"finished_at"`
	Result        string         `json:"result"`
	Cleanup       string         `json:"cleanup"`
	Events        []event        `json:"events"`
	Environment   runmeta.Report `json:"environment"`
}

func Run(ctx context.Context, out io.Writer) (finalErr error) {
	return runHost(ctx, out, false)
}

// RunInfrastructure verifies PostgreSQL, Kafka and Redis together.
func RunInfrastructure(ctx context.Context, out io.Writer) error {
	return runHost(ctx, out, true)
}

func runHost(ctx context.Context, out io.Writer, fullStack bool) error {
	root, err := os.Getwd()
	if err != nil {
		return errors.New("cannot locate checkout")
	}
	return runMode(ctx, out, root, nil, fullStack)
}

func run(ctx context.Context, out io.Writer, root string, execute func(context.Context, string, []string) (command.Output, error)) (finalErr error) {
	return runMode(ctx, out, root, execute, false)
}

func runMode(ctx context.Context, out io.Writer, root string, execute func(context.Context, string, []string) (command.Output, error), fullStack bool) (finalErr error) {
	var err error
	var token [8]byte
	if _, err = rand.Read(token[:]); err != nil {
		return errors.New("cannot generate project ID")
	}
	project := "faultline-data-test-" + hex.EncodeToString(token[:])
	artifactKind := "postgres"
	if fullStack {
		artifactKind = "infrastructure"
	}
	dir := filepath.Join(root, "artifacts", artifactKind, project)
	if err = os.MkdirAll(dir, 0o700); err != nil {
		return errors.New("cannot create report directory")
	}
	cfg, err := database.Load(filepath.Join(root, "configs/postgres/local.json"), filepath.Join(dir, "secrets"))
	if err != nil {
		return err
	}
	var infraCfg infrastructure.Config
	if fullStack {
		infraCfg, err = infrastructure.Load(filepath.Join(root, "configs/infrastructure/local.json"), filepath.Join(dir, "infrastructure-secrets"))
		if err != nil {
			return err
		}
	}
	compose := func(args ...string) []string { return composeModeArgs(project, fullStack, args...) }
	removeCredentials := func() error {
		var result error
		// Windows/OneDrive can mark a bind-source directory read-only.
		// Restore write access only to this run's owned credential directory.
		if e := os.Chmod(cfg.Secrets, 0o700); e != nil && !errors.Is(e, os.ErrNotExist) {
			result = errors.Join(result, errors.New("cannot restore temporary directory permissions"))
		}
		for _, role := range append([]string{cfg.AdminUser}, roleNames(cfg)...) {
			if e := os.Remove(filepath.Join(cfg.Secrets, role+".password")); e != nil && !errors.Is(e, os.ErrNotExist) {
				result = errors.Join(result, errors.New("cannot remove temporary credential"))
			}
		}
		if e := os.Remove(cfg.Secrets); e != nil && !errors.Is(e, os.ErrNotExist) {
			result = errors.Join(result, errors.New("cannot remove temporary credential directory"))
		}
		if fullStack {
			if e := os.Chmod(infraCfg.Secrets, 0o700); e != nil && !errors.Is(e, os.ErrNotExist) {
				result = errors.Join(result, errors.New("cannot restore infrastructure credential directory permissions"))
			}
			for _, name := range infrastructure.CredentialFiles {
				if e := os.Remove(filepath.Join(infraCfg.Secrets, name)); e != nil && !errors.Is(e, os.ErrNotExist) {
					result = errors.Join(result, errors.New("cannot remove infrastructure credential"))
				}
			}
			if e := os.Remove(infraCfg.Secrets); e != nil && !errors.Is(e, os.ErrNotExist) {
				result = errors.Join(result, errors.New("cannot remove infrastructure credential directory"))
			}
		}
		return result
	}
	defer func() { finalErr = errors.Join(finalErr, removeCredentials()) }()
	if err = database.GenerateSecrets(cfg); err != nil {
		return err
	}
	if fullStack {
		if err = infrastructure.GenerateSecrets(infraCfg); err != nil {
			return err
		}
	}
	pairs := []string{root, "[workspace]", filepath.ToSlash(root), "[workspace]"}
	if home, e := os.UserHomeDir(); e == nil && home != "" {
		pairs = append(pairs, home, "[home]", filepath.ToSlash(home), "[home]")
	}
	for _, role := range append([]string{cfg.AdminUser}, roleNames(cfg)...) {
		value, e := cfg.Secret(role)
		if e != nil {
			return e
		}
		pairs = append(pairs, value, "[secret]")
	}
	if fullStack {
		value, e := infrastructure.Password(infraCfg)
		if e != nil {
			return e
		}
		pairs = append(pairs, value, "[secret]", infrastructure.PasswordHash(value), "[secret hash]")
	}
	redact := strings.NewReplacer(pairs...)
	env := []string{}
	for _, item := range os.Environ() {
		key, _, _ := strings.Cut(item, "=")
		if !strings.EqualFold(key, "FL_PG_SECRETS_DIR") && !strings.EqualFold(key, "FL_INFRA_SECRETS_DIR") && !strings.EqualFold(key, "FL_INVENTORY_DATASET_ID") {
			env = append(env, item)
		}
	}
	env = append(env, "FL_PG_SECRETS_DIR="+cfg.Secrets)
	// Each disposable database gets a unique namespace, independent of test keys.
	env = append(env, "FL_INVENTORY_DATASET_ID="+project)
	if fullStack {
		env = append(env, "FL_INFRA_SECRETS_DIR="+infraCfg.Secrets)
	}
	executor := command.Executor{Directory: root, Env: env}
	if execute == nil {
		execute = executor.Execute
	}
	services := []string{"postgres"}
	if fullStack {
		services = append(services, "kafka", "redis")
	}
	observedServices := append(append([]string{}, services...), "inventory")
	evidence := report{SchemaVersion: 1, Project: project, Started: time.Now().UTC(), Result: "FAIL", Cleanup: "PENDING", Events: []event{}, Environment: runmeta.New(observedServices)}
	save := func() error {
		data, e := json.MarshalIndent(evidence, "", "  ")
		if e != nil {
			return e
		}
		return os.WriteFile(filepath.Join(dir, "report.json"), append(data, '\n'), 0o600)
	}
	if err = save(); err != nil {
		return errors.New("cannot save initial report")
	}
	call := func(parent context.Context, stage string, budget time.Duration, args ...string) (string, error) {
		_, _ = fmt.Fprintln(out, artifactKind+" smoke:", stage)
		commandCtx, cancel := context.WithTimeout(parent, budget)
		defer cancel()
		result, e := execute(commandCtx, "docker", args)
		observation := event{Stage: stage, Output: redact.Replace(result.Stdout + result.Stderr), Passed: e == nil, ExitCode: result.ExitCode}
		if e != nil {
			observation.Error = redact.Replace(e.Error())
			if commandCtx.Err() != nil {
				observation.Error = redact.Replace(errors.Join(commandCtx.Err(), e).Error())
			}
			e = fmt.Errorf("%s failed: %s; see redacted report", stage, observation.Error)
		}
		evidence.Events = append(evidence.Events, observation)
		return result.Stdout, errors.Join(e, save())
	}
	mayExist := false
	defer func() {
		if mayExist {
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			_, e := call(cleanupCtx, "cleanup", 60*time.Second, compose("down", "--volumes", "--remove-orphans", "--rmi", "local")...)
			if e == nil {
				for _, args := range [][]string{
					{"ps", "-aq", "--filter", "label=com.docker.compose.project=" + project},
					{"network", "ls", "-q", "--filter", "label=com.docker.compose.project=" + project},
					{"volume", "ls", "-q", "--filter", "label=com.docker.compose.project=" + project},
					{"image", "ls", "-q", "--filter", "reference=" + project + "-*:latest"},
				} {
					var remaining string
					remaining, e = call(cleanupCtx, "verify cleanup", 5*time.Second, args...)
					if e != nil || strings.TrimSpace(remaining) != "" {
						e = errors.New("owned resources remain after cleanup")
						break
					}
				}
			}
			evidence.Cleanup = "PASS"
			if e != nil {
				evidence.Cleanup = "FAIL"
			}
			finalErr = errors.Join(finalErr, e)
		} else {
			evidence.Cleanup = "NOT_NEEDED"
		}
		credentialErr := removeCredentials()
		if credentialErr != nil {
			evidence.Cleanup = "FAIL"
		}
		finalErr = errors.Join(finalErr, credentialErr)
		evidence.Finished = time.Now().UTC()
		if finalErr == nil {
			evidence.Result = "PASS"
		}
		finalErr = errors.Join(finalErr, save())
		_, _ = fmt.Fprintf(out, "%s: artifacts/%s/%s/report.json\n", evidence.Result, artifactKind, project)
	}()
	if _, err = call(ctx, "prerequisites", 15*time.Second, "version", "--format", "{{.Server.Version}}"); err != nil {
		return err
	}
	mayExist = true
	buildServices := []string{"dbtool"}
	buildServices = append(buildServices, "inventory")
	if fullStack {
		buildServices = append(buildServices, "infrastructure-check")
	}
	// Register before build/start, so partial startup failures retain diagnostics.
	defer func() {
		logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		logServices := append(append([]string{}, services...), "inventory")
		_, _ = call(logCtx, "dependency logs", 5*time.Second, compose(append([]string{"logs", "--no-color", "--tail", "200"}, logServices...)...)...)
	}()
	if _, err = call(ctx, "build dbtool", 10*time.Minute, compose(append([]string{"build"}, buildServices...)...)...); err != nil {
		return err
	}
	_, err = call(ctx, "start dependencies", 2*time.Minute, compose(append([]string{"up", "-d", "--wait", "--wait-timeout", "100"}, services...)...)...)
	if err != nil {
		return err
	}
	tool := func(stage string, args ...string) (string, error) {
		return call(ctx, stage, 90*time.Second, compose(append([]string{"run", "--rm", "--no-deps", "dbtool"}, args...)...)...)
	}
	for attempt := 0; attempt < 2; attempt++ {
		if _, err = tool("bootstrap", "bootstrap"); err != nil {
			return err
		}
	}
	if _, err = call(ctx, "rights, migrations and seed", 3*time.Minute, compose("run", "--rm", "--no-deps", "-e", "FL_DATA_INTEGRATION=1", "--entrypoint", "/app/database.test", "dbtool", "-test.run=^TestPostgresFoundation$", "-test.v", "-test.timeout=150s")...); err != nil {
		return err
	}
	if _, inventoryErr := call(ctx, "Inventory operations and process restart", 3*time.Minute, compose("run", "--rm", "--no-deps", "-e", "FL_DATA_INTEGRATION=1", "--entrypoint", "/app/inventory.test", "dbtool", "-test.run=^TestInventory(Foundation|Process)$", "-test.v", "-test.timeout=150s")...); inventoryErr != nil {
		return inventoryErr
	}
	// Put the business overlay last, preserving its command/config/secret wiring.
	base := compose()
	var deployed []string
	for index, value := range base {
		if value == "-p" {
			deployed = append(append(append([]string{}, base[:index]...), "-f", "compose.inventory.yaml"), base[index:]...)
			break
		}
	}
	_, err = call(ctx, "start Inventory business container", 60*time.Second, append(deployed, "up", "-d", "--no-deps", "--wait", "--wait-timeout", "45", "inventory")...)
	// Observe the business container with its effective overlay, before its
	// functional probes. Raw config bypasses call() to keep secrets out of logs.
	evidence.Environment = runmeta.Collect(ctx, func(ctx context.Context, args []string) (string, error) {
		result, e := execute(ctx, "docker", args)
		if e == nil && result.ExitCode != 0 {
			e = errors.New("metadata command failed")
		}
		return result.Stdout, e
	}, deployed, observedServices)
	if err = errors.Join(err, ctx.Err(), save()); err != nil {
		return err
	}
	deployedProbe := func(stage string) error {
		_, e := call(ctx, stage, 45*time.Second, compose("run", "--rm", "--no-deps", "-e", "FL_DATA_INTEGRATION=1", "-e", "FL_INVENTORY_TARGET=inventory:50051", "--entrypoint", "/app/inventory.test", "dbtool", "-test.run=^TestInventoryDeployed$", "-test.v", "-test.timeout=30s")...)
		return e
	}
	if err = deployedProbe("Inventory deployed gRPC API"); err != nil {
		return err
	}
	// Inventory tests add owned fixtures; capture the baseline after their writes.
	before, err := tool("checksum before restart", "checksum")
	if err != nil {
		return err
	}
	var checkpoint string
	if fullStack {
		checkpoint, err = call(ctx, "Kafka and Redis before restart", 60*time.Second, compose("run", "--rm", "--no-deps", "infrastructure-check", "check-infrastructure", "-stage", "before-restart")...)
		if err != nil {
			return err
		}
		var plan infrastructure.Checkpoint
		if err = json.Unmarshal([]byte(checkpoint), &plan); err != nil || plan.Version != 1 {
			return errors.New("infrastructure checker returned an invalid checkpoint")
		}
	}
	if _, err = call(ctx, "restart dependencies", 90*time.Second, compose(append([]string{"restart"}, services...)...)...); err != nil {
		return err
	}
	if _, err = call(ctx, "wait after restart", 2*time.Minute, compose(append([]string{"up", "-d", "--wait", "--wait-timeout", "100"}, services...)...)...); err != nil {
		return err
	}
	for _, operation := range []string{"bootstrap", "migrate", "seed"} {
		if _, err = tool("repeat "+operation, operation); err != nil {
			return err
		}
	}
	after, err := tool("checksum after restart", "checksum")
	if err != nil {
		return err
	}
	if strings.TrimSpace(before) != strings.TrimSpace(after) {
		return errors.New("data or migration history changed across restart/repeat")
	}
	if _, err = call(ctx, "Inventory persistence after PostgreSQL restart", time.Minute, compose("run", "--rm", "--no-deps", "-e", "FL_DATA_INTEGRATION=1", "--entrypoint", "/app/inventory.test", "dbtool", "-test.run=^TestInventoryAfterRestart$", "-test.v", "-test.timeout=45s")...); err != nil {
		return err
	}
	if err = deployedProbe("Inventory deployed API after PostgreSQL restart"); err != nil {
		return err
	}
	if fullStack {
		if _, err = call(ctx, "Kafka and Redis after restart", 60*time.Second, compose("run", "--rm", "--no-deps", "infrastructure-check", "check-infrastructure", "-stage", "after-restart", "-checkpoint", strings.TrimSpace(checkpoint))...); err != nil {
			return err
		}
	}
	return nil
}
func roleNames(cfg database.Config) []string {
	var result []string
	for _, service := range cfg.Services {
		result = append(result, service.AppUser, service.MigrationUser)
	}
	return result
}
