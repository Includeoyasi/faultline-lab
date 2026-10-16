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
)

var projectName = regexp.MustCompile(`^faultline-data-[a-z0-9][a-z0-9-]{0,47}$`)

func composeArgs(project string, args ...string) []string {
	return append([]string{"compose", "-f", "compose.yaml", "-f", "compose.postgres.yaml", "-p", project}, args...)
}

// Reset requires an exact explicit project confirmation; ordinary down preserves data.
func Reset(ctx context.Context, project, confirm string, out io.Writer) error {
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
	result, err := executor.Execute(bounded, "docker", composeArgs(project, "down", "--volumes", "--remove-orphans", "--rmi", "local"))
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
	Project  string    `json:"project"`
	Started  time.Time `json:"started_at"`
	Finished time.Time `json:"finished_at"`
	Result   string    `json:"result"`
	Cleanup  string    `json:"cleanup"`
	Events   []event   `json:"events"`
}

func Run(ctx context.Context, out io.Writer) (finalErr error) {
	root, err := os.Getwd()
	if err != nil {
		return errors.New("cannot locate checkout")
	}
	return run(ctx, out, root, nil)
}

func run(ctx context.Context, out io.Writer, root string, execute func(context.Context, string, []string) (command.Output, error)) (finalErr error) {
	var err error
	var token [8]byte
	if _, err = rand.Read(token[:]); err != nil {
		return errors.New("cannot generate project ID")
	}
	project := "faultline-data-test-" + hex.EncodeToString(token[:])
	dir := filepath.Join(root, "artifacts", "postgres", project)
	if err = os.MkdirAll(dir, 0o700); err != nil {
		return errors.New("cannot create report directory")
	}
	cfg, err := database.Load(filepath.Join(root, "configs/postgres/local.json"), filepath.Join(dir, "secrets"))
	if err != nil {
		return err
	}
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
		return result
	}
	defer func() { finalErr = errors.Join(finalErr, removeCredentials()) }()
	if err = database.GenerateSecrets(cfg); err != nil {
		return err
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
	redact := strings.NewReplacer(pairs...)
	env := []string{}
	for _, item := range os.Environ() {
		key, _, _ := strings.Cut(item, "=")
		if !strings.EqualFold(key, "FL_PG_SECRETS_DIR") {
			env = append(env, item)
		}
	}
	env = append(env, "FL_PG_SECRETS_DIR="+cfg.Secrets)
	executor := command.Executor{Directory: root, Env: env}
	if execute == nil {
		execute = executor.Execute
	}
	evidence := report{Project: project, Started: time.Now().UTC(), Result: "FAIL", Cleanup: "PENDING", Events: []event{}}
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
		_, _ = fmt.Fprintln(out, "PostgreSQL smoke:", stage)
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
			_, e := call(cleanupCtx, "cleanup", 60*time.Second, composeArgs(project, "down", "--volumes", "--remove-orphans", "--rmi", "local")...)
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
		_, _ = fmt.Fprintf(out, "%s: artifacts/postgres/%s/report.json\n", evidence.Result, project)
	}()
	if _, err = call(ctx, "prerequisites", 15*time.Second, "version", "--format", "{{.Server.Version}}"); err != nil {
		return err
	}
	mayExist = true
	if _, err = call(ctx, "build dbtool", 10*time.Minute, composeArgs(project, "build", "dbtool")...); err != nil {
		return err
	}
	if _, err = call(ctx, "start PostgreSQL", 90*time.Second, composeArgs(project, "up", "-d", "--wait", "--wait-timeout", "75", "postgres")...); err != nil {
		return err
	}
	defer func() {
		logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = call(logCtx, "PostgreSQL logs", 5*time.Second, composeArgs(project, "logs", "--no-color", "postgres")...)
	}()
	tool := func(stage string, args ...string) (string, error) {
		return call(ctx, stage, 90*time.Second, composeArgs(project, append([]string{"run", "--rm", "--no-deps", "dbtool"}, args...)...)...)
	}
	for attempt := 0; attempt < 2; attempt++ {
		if _, err = tool("bootstrap", "bootstrap"); err != nil {
			return err
		}
	}
	if _, err = call(ctx, "rights, migrations and seed", 3*time.Minute, composeArgs(project, "run", "--rm", "--no-deps", "-e", "FL_DATA_INTEGRATION=1", "--entrypoint", "/app/database.test", "dbtool", "-test.run=^TestPostgresFoundation$", "-test.v", "-test.timeout=150s")...); err != nil {
		return err
	}
	before, err := tool("fingerprint before restart", "fingerprint")
	if err != nil {
		return err
	}
	if _, err = call(ctx, "restart PostgreSQL", 60*time.Second, composeArgs(project, "restart", "postgres")...); err != nil {
		return err
	}
	if _, err = call(ctx, "wait after restart", 60*time.Second, composeArgs(project, "up", "-d", "--wait", "--wait-timeout", "45", "postgres")...); err != nil {
		return err
	}
	for _, operation := range []string{"bootstrap", "migrate", "seed"} {
		if _, err = tool("repeat "+operation, operation); err != nil {
			return err
		}
	}
	after, err := tool("fingerprint after restart", "fingerprint")
	if err != nil {
		return err
	}
	if strings.TrimSpace(before) != strings.TrimSpace(after) {
		return errors.New("data or migration history changed across restart/repeat")
	}
	_, err = call(ctx, "resource snapshot", 15*time.Second, "stats", "--no-stream", "--format", "{{json .}}", project+"-postgres-1")
	return err
}
func roleNames(cfg database.Config) []string {
	var result []string
	for _, service := range cfg.Services {
		result = append(result, service.AppUser, service.MigrationUser)
	}
	return result
}
