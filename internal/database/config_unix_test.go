//go:build linux || darwin

package database

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestSecretsReadableByContainerDespiteRestrictiveUmask(t *testing.T) {
	previous := syscall.Umask(0o077)
	defer syscall.Umask(previous)
	cfg := testConfig(t)
	if err := GenerateSecrets(cfg); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(cfg.Secrets, cfg.AdminUser+".password")
	value, err := cfg.Secret(cfg.AdminUser)
	if err != nil {
		t.Fatal(err)
	}
	for _, mode := range []os.FileMode{0o644, 0o600} {
		if mode == 0o600 {
			if err = os.Chmod(path, mode); err != nil {
				t.Fatal(err)
			}
			if err = GenerateSecrets(cfg); err != nil {
				t.Fatal(err)
			}
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0o644 {
			t.Fatalf("container cannot read secret: mode=%o", info.Mode().Perm())
		}
		after, err := cfg.Secret(cfg.AdminUser)
		if err != nil || after != value {
			t.Fatal("credentials changed", err)
		}
	}
}
