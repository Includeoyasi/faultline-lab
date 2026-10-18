package database

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
)

// Checksum covers ordered migration history and Inventory product-stock pairs.
// It is evidence of persistence for this foundation, not a general backup tool.
func Checksum(ctx context.Context, cfg Config) (string, error) {
	var values []string
	for _, service := range cfg.Services {
		conn, err := Connect(ctx, cfg, service.Database, service.MigrationUser)
		if err != nil {
			return "", err
		}
		var history string
		err = conn.QueryRow(ctx, `SELECT coalesce(json_agg(v ORDER BY version)::text,'[]') FROM lab_migrations.versions v`).Scan(&history)
		Close(conn)
		if err != nil {
			return "", SafeError("history checksum", err)
		}
		values = append(values, service.Name, history)
	}
	service, err := cfg.Service("inventory")
	if err != nil {
		return "", err
	}
	conn, err := Connect(ctx, cfg, service.Database, service.AppUser)
	if err != nil {
		return "", err
	}
	defer Close(conn)
	var catalog string
	err = conn.QueryRow(ctx, `SELECT coalesce(json_agg(v ORDER BY sku)::text,'[]') FROM (SELECT p.*,s.available FROM inventory.products p JOIN inventory.stock s USING(sku)) v`).Scan(&catalog)
	if err != nil {
		return "", SafeError("catalog checksum", err)
	}
	values = append(values, catalog)
	data, err := json.Marshal(values)
	if err != nil {
		return "", errors.New("cannot encode checksum")
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
