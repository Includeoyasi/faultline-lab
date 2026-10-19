package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestOperationErrorsExitNonzero(t *testing.T) {
	var out, errs bytes.Buffer
	code := run(context.Background(), []string{"migrate", "-config", "../../configs/postgres/local.json", "-secrets", t.TempDir(), "-service", "unknown"}, &out, &errs)
	if code != 1 || !strings.Contains(errs.String(), "unknown service") || strings.Contains(out.String(), "PASS") {
		t.Fatalf("code=%d out=%q err=%q", code, out.String(), errs.String())
	}
}

func TestUsageAndSecrets(t *testing.T) {
	for _, test := range []struct {
		args []string
		code int
	}{
		{nil, 2}, {[]string{"unknown"}, 2}, {[]string{"secrets", "-help"}, 0},
		{[]string{"secrets", "-config", "../../configs/postgres/local.json", "-secrets", t.TempDir()}, 0},
	} {
		var out, errs bytes.Buffer
		if code := run(context.Background(), test.args, &out, &errs); code != test.code {
			t.Fatalf("args=%v code=%d err=%s", test.args, code, errs.String())
		}
	}
}

func TestInfrastructureCheckDoesNotLoadPostgresCredentials(t *testing.T) {
	var out, errs bytes.Buffer
	code := run(t.Context(), []string{"check-infrastructure", "-config", "missing-postgres.json", "-infrastructure-config", "../../configs/infrastructure/local.json", "-infrastructure-secrets", t.TempDir(), "-stage", "invalid"}, &out, &errs)
	if code != 1 || !strings.Contains(errs.String(), "stage must be") || strings.Contains(errs.String(), "PostgreSQL") || out.Len() != 0 {
		t.Fatalf("code=%d out=%q err=%q", code, out.String(), errs.String())
	}
}
