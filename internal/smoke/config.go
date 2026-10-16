package smoke

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// ShutdownMargin reserves time outside the application's RPC drain deadline.
const ShutdownMargin = 2 * time.Second

// CheckConfig checks the user's effective deployment configuration, including
// the selected env file. Unlike Run it does not force the smoke fixture profile.
// Resolved config can contain secrets: it stays in memory and is never logged.
func CheckConfig(ctx context.Context, root string, allowUnsafe bool, output io.Writer) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	executor := processExecutor{directory: root, env: os.Environ()}
	resolved, err := executor.execute(ctx, "docker", []string{"compose", "-f", "compose.yaml", "config", "--format", "json"})
	if err != nil {
		return fmt.Errorf("cannot resolve Compose config (exit %d); check Docker CLI and compose.yaml", resolved.exitCode)
	}
	warnings, err := validateConfig([]byte(resolved.stdout), allowUnsafe)
	if err != nil {
		return err
	}
	for _, warning := range warnings {
		_, _ = fmt.Fprintln(output, "WARNING:", warning)
	}
	_, err = fmt.Fprintln(output, "Compose shutdown budgets checked; no services started")
	return err
}

func validateConfig(data []byte, allowUnsafe bool) ([]string, error) {
	var config struct {
		Services map[string]struct {
			Environment map[string]string `json:"environment"`
			Grace       string            `json:"stop_grace_period"`
		} `json:"services"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid resolved Compose JSON: %w", err)
	}
	var warnings []string
	for _, name := range services {
		service, ok := config.Services[name]
		if !ok {
			return nil, fmt.Errorf("compose config lacks service %s", name)
		}
		appValue, present := service.Environment["FL_SHUTDOWN_TIMEOUT"]
		if !present {
			appValue = "10s" // The application's documented default.
		}
		app, err := time.ParseDuration(appValue)
		if err != nil || app <= 0 {
			return nil, fmt.Errorf("%s: FL_SHUTDOWN_TIMEOUT must be a positive duration", name)
		}
		graceValue := service.Grace
		if graceValue == "" {
			graceValue = "10s" // Compose's default when the property is omitted.
		}
		grace, err := time.ParseDuration(graceValue)
		if err != nil || grace <= 0 {
			return nil, fmt.Errorf("%s: stop_grace_period must be a positive duration", name)
		}
		// Subtraction avoids overflow for an otherwise valid very long duration.
		if grace < app || grace-app < ShutdownMargin {
			message := fmt.Sprintf("%s: container grace %s must cover app timeout %s plus %s margin", name, grace, app, ShutdownMargin)
			if !allowUnsafe {
				return nil, fmt.Errorf("%s", message)
			}
			warnings = append(warnings, message)
		}
	}
	return warnings, nil
}
