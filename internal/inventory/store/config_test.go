package store

import (
	"testing"
)

func TestConfigRestrictsIdentityAndBudgets(t *testing.T) {
	base := Config{Host: "postgres", Port: 5432, Database: "inventory_db", User: "inventory_app", SSLMode: "disable", PasswordFile: "/run/secrets/inventory_app.password", DatasetID: "dataset", ConnectTimeout: "5s", RequestTimeout: "3s", MaxConnections: 4}
	if err := base.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*Config){func(c *Config) { c.User = "postgres" }, func(c *Config) { c.User = "inventory_migrator" }, func(c *Config) { c.Database = "orders_db" }, func(c *Config) { c.MaxConnections = 30 }, func(c *Config) { c.RequestTimeout = "0s" }, func(c *Config) { c.DatasetID = "" }} {
		cfg := base
		mutate(&cfg)
		if cfg.Validate() == nil {
			t.Fatalf("invalid config accepted: %+v", cfg)
		}
	}
}
