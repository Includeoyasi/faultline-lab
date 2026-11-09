package orders

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigurationBudgetsAndIdentity(t *testing.T) {
	c, err := LoadConfig("../../configs/orders/local.json", "dataset")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"foreign-user", "no-peer", "bad-peer", "unbounded-workers", "short-lease", "zero-retry", "reversed-retry"} {
		t.Run(name, func(t *testing.T) {
			invalid := c
			switch name {
			case "foreign-user":
				invalid.Database.User = "inventory_app"
			case "no-peer":
				invalid.InventoryTarget = ""
			case "bad-peer":
				invalid.InventoryTarget = "inventory:0"
			case "unbounded-workers":
				invalid.WorkerCount = 100
			case "short-lease":
				invalid.Lease = "4s"
			case "zero-retry":
				invalid.RetryMin = "0s"
			case "reversed-retry":
				invalid.RetryMin = "10s"
			}
			if err = invalid.Validate(); err == nil {
				t.Fatal("invalid configuration accepted")
			}
		})
	}
	data, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "config.json")
	if err = os.WriteFile(path, append(data, []byte(` {"trailing":true}`)...), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err = LoadConfig(path, "dataset"); err == nil {
		t.Fatal("trailing configuration accepted")
	}
}
