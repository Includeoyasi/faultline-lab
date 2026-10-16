package smoke

import (
	"encoding/json"
	"testing"
)

func resolvedConfig(app, grace string) []byte {
	config := map[string]any{"services": map[string]any{}}
	for _, name := range services {
		config["services"].(map[string]any)[name] = map[string]any{"environment": map[string]string{"FL_SHUTDOWN_TIMEOUT": app}, "stop_grace_period": grace}
	}
	data, _ := json.Marshal(config)
	return data
}

func TestShutdownBudgets(t *testing.T) {
	for _, tc := range []struct {
		name, app, grace string
		allow, valid     bool
		warnings         int
	}{
		{"defaults", "10s", "15s", false, true, 0},
		{"minimum-margin", "10s", "12s", false, true, 0},
		{"external-too-short", "20s", "15s", false, false, 0},
		{"no-margin", "15s", "15s", false, false, 0},
		{"explicit-experiment", "20s", "15s", true, true, 3},
		{"invalid-app", "secret-value", "15s", true, false, 0},
		{"invalid-grace", "10s", "invalid", true, false, 0},
		{"zero-app", "0s", "15s", false, false, 0},
		{"negative-grace", "10s", "-1s", false, false, 0},
		{"overflow-margin", "2562047h47m16.854775807s", "2562047h47m16.854775807s", false, false, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			warnings, err := validateConfig(resolvedConfig(tc.app, tc.grace), tc.allow)
			if (err == nil) != tc.valid || len(warnings) != tc.warnings {
				t.Fatalf("warnings=%v error=%v", warnings, err)
			}
		})
	}
	for _, data := range []string{`not-json`, `{"services":{}}`} {
		if _, err := validateConfig([]byte(data), false); err == nil {
			t.Fatal("invalid/missing services accepted")
		}
	}
}
