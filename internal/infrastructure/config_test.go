package infrastructure

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testConfig(t *testing.T) Config {
	t.Helper()
	cfg, err := Load("../../configs/infrastructure/local.json", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func TestConfigRejectsUnsafeProfiles(t *testing.T) {
	for _, change := range []func(*Config){
		func(c *Config) { c.Timeout = "0s" },
		func(c *Config) { c.Timeout = "31s" },
		func(c *Config) { c.ProbeTTL = c.Timeout },
		func(c *Config) { c.KafkaBrokers = nil },
		func(c *Config) { c.RedisAddress = "redis" },
		func(c *Config) { c.RedisUser = "lab on nopass" },
		func(c *Config) { c.RedisUser = "default" },
		func(c *Config) { c.KeyPrefix = "*" },
	} {
		cfg := testConfig(t)
		change(&cfg)
		if cfg.Validate() == nil {
			t.Fatal("accepted unsafe profile")
		}
	}
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"unknown":true}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path, t.TempDir()); err == nil {
		t.Fatal("accepted unknown configuration")
	}
}

func TestCredentialsPreservedAndACLRestricted(t *testing.T) {
	cfg := testConfig(t)
	if err := GenerateSecrets(cfg); err != nil {
		t.Fatal(err)
	}
	password, err := Password(cfg)
	if err != nil {
		t.Fatal(err)
	}
	acl, err := os.ReadFile(filepath.Join(cfg.Secrets, "redis.acl"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(acl), password) || !strings.Contains(string(acl), "#"+PasswordHash(password)) || !strings.Contains(string(acl), "user default off") || !strings.Contains(string(acl), "~"+cfg.KeyPrefix+"* -@all") {
		t.Fatal("ACL exposes plaintext or does not restrict commands/keys")
	}
	if err = GenerateSecrets(cfg); err != nil {
		t.Fatal(err)
	}
	again, _ := Password(cfg)
	if again != password {
		t.Fatal("repeat changed password")
	}
	if err = os.WriteFile(filepath.Join(cfg.Secrets, "redis.acl"), []byte("unexpected configuration"), 0o600); err != nil {
		t.Fatal(err)
	}
	if GenerateSecrets(cfg) == nil {
		t.Fatal("silently replaced an inconsistent ACL")
	}
}

func TestInvalidCheckpointFailsBeforeConnections(t *testing.T) {
	cfg := testConfig(t)
	for _, saved := range []string{"", `{}`, `{"version":1,"topic":"foreign"}`, `{"unknown":true}`, strings.Repeat("x", 2049)} {
		if _, err := Check(t.Context(), cfg, "after-restart", saved); err == nil || strings.Contains(err.Error(), "credential") {
			t.Fatal("invalid checkpoint reached credential loading")
		}
	}
}
