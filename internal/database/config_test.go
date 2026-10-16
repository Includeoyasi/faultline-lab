package database

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/jackc/pgx/v5/pgconn"
)

func testConfig(t *testing.T) Config {
	t.Helper()
	cfg, err := Load("../../configs/postgres/local.json", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}
func TestConfigValidation(t *testing.T) {
	for _, timeout := range []string{"0", "-1s", "500us", "61s", "invalid"} {
		t.Run(timeout, func(t *testing.T) {
			cfg := testConfig(t)
			cfg.Timeout = timeout
			if cfg.Validate() == nil {
				t.Fatal("accepted invalid budget")
			}
		})
	}
	cfg := testConfig(t)
	cfg.Services[0].AppUser = cfg.Services[1].AppUser
	if cfg.Validate() == nil {
		t.Fatal("accepted shared role")
	}
	cfg = testConfig(t)
	cfg.Services[0].Database = "db; DROP"
	if cfg.Validate() == nil {
		t.Fatal("accepted SQL identifier")
	}
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"unknown":true}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path, t.TempDir()); err == nil {
		t.Fatal("accepted unknown settings")
	}
}
func TestSecretsPreservedAndMalformedRejected(t *testing.T) {
	cfg := testConfig(t)
	if err := GenerateSecrets(cfg); err != nil {
		t.Fatal(err)
	}
	first, err := cfg.Secret(cfg.AdminUser)
	if err != nil {
		t.Fatal(err)
	}
	if err = GenerateSecrets(cfg); err != nil {
		t.Fatal(err)
	}
	second, _ := cfg.Secret(cfg.AdminUser)
	if first != second {
		t.Fatal("credentials rotated unexpectedly")
	}
	if err = os.WriteFile(filepath.Join(cfg.Secrets, cfg.AdminUser+".password"), []byte("malformed-sensitive-value"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err = GenerateSecrets(cfg); err == nil || strings.Contains(err.Error(), "malformed-sensitive-value") {
		t.Fatal("malformed secret accepted or exposed")
	}
}
func TestSafeError(t *testing.T) {
	for _, err := range []error{&pgconn.PgError{Code: "42501", Message: "secret DSN", Detail: "password"}, errors.New("password"), errors.Join(context.Canceled, errors.New("password"))} {
		safe := SafeError("operation", err)
		if strings.Contains(safe.Error(), "password") || strings.Contains(safe.Error(), "DSN") {
			t.Fatal("sensitive diagnostics leaked")
		}
	}
	if !errors.Is(SafeError("operation", context.DeadlineExceeded), context.DeadlineExceeded) {
		t.Fatal("deadline lost")
	}
}
func TestMigrationFiles(t *testing.T) {
	for _, files := range []fstest.MapFS{
		{"orders/0000_invalid.sql": {Data: []byte("SELECT 1")}},
		{"orders/0001_a.sql": {Data: []byte("SELECT 1")}, "orders/0001_b.sql": {Data: []byte("SELECT 2")}},
		{"orders/unversioned.sql": {Data: []byte("SELECT 1")}},
		{"orders/0001_empty.sql": {Data: []byte{}}},
	} {
		if _, err := loadMigrations(files, "orders"); err == nil {
			t.Fatal("invalid migration set accepted")
		}
	}
	files := fstest.MapFS{"orders/0001_a.sql": {Data: []byte("SELECT 1")}}
	before, err := loadMigrations(files, "orders")
	if err != nil {
		t.Fatal(err)
	}
	files["orders/0001_a.sql"].Data = []byte("SELECT 2")
	after, err := loadMigrations(files, "orders")
	if err != nil {
		t.Fatal(err)
	}
	if before[0].checksum == after[0].checksum {
		t.Fatal("SQL edit did not alter checksum")
	}
}
