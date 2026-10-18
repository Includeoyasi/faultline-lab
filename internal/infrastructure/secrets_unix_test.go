//go:build !windows

package infrastructure

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestSecretPermissionsUnderRestrictiveUmask(t *testing.T) {
	previous := syscall.Umask(0o077)
	defer syscall.Umask(previous)
	cfg := testConfig(t)
	if err := GenerateSecrets(cfg); err != nil {
		t.Fatal(err)
	}
	for _, name := range CredentialFiles {
		info, err := os.Stat(filepath.Join(cfg.Secrets, name))
		if err != nil || info.Mode().Perm() != 0o644 {
			t.Fatal("credential cannot be consumed by its non-root container")
		}
	}
}
