package inventory

import (
	"errors"
	"testing"
)

func TestCanonicalOperationIdentity(t *testing.T) {
	base := Request{DatasetID: "dataset", OrderID: "order", OperationID: "key", Method: "reserve", Items: []Item{{"B", 2}, {"A", 1}}}
	normalized, hash, err := canonical(base)
	if err != nil {
		t.Fatal(err)
	}
	if base.Items[0].SKU != "B" || normalized.Items[0].SKU != "A" {
		t.Fatal("canonicalization mutates caller or does not sort")
	}
	reordered := base
	reordered.Items = []Item{{"A", 1}, {"B", 2}}
	if _, same, e := canonical(reordered); e != nil || hash != same {
		t.Fatal("SKU order changed operation identity", e)
	}
	changed := base
	changed.OrderID = "other"
	if _, different, e := canonical(changed); e != nil || hash == different {
		t.Fatal("different order reused request identity", e)
	}
}

func TestInvalidRequests(t *testing.T) {
	for _, r := range []Request{
		{},
		{DatasetID: "d", OrderID: "o", OperationID: "k", Method: "reserve"},
		{DatasetID: "d", OrderID: "o", OperationID: "k", Method: "reserve", Items: []Item{{"A", 0}}},
		{DatasetID: "d", OrderID: "o", OperationID: "k", Method: "reserve", Items: []Item{{"A", 1}, {"A", 2}}},
		{DatasetID: "d", OrderID: "o", OperationID: "k", Method: "commit", Items: []Item{{"A", 1}}},
		{DatasetID: "../secret", OrderID: "o", OperationID: "k", Method: "release"},
	} {
		if _, _, err := canonical(r); !errors.Is(err, ErrInvalid) {
			t.Fatalf("request accepted: %+v", r)
		}
	}
}

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
