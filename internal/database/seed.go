package database

import "context"

// Seed inserts missing Inventory fixtures without overwriting existing values.
func Seed(ctx context.Context, cfg Config) error {
	service, err := cfg.Service("inventory")
	if err != nil {
		return err
	}
	conn, err := Connect(ctx, cfg, service.Database, service.AppUser)
	if err != nil {
		return err
	}
	defer Close(conn)
	tx, err := conn.Begin(ctx)
	if err != nil {
		return SafeError("seed transaction", err)
	}
	defer rollback(tx)
	if _, err = tx.Exec(ctx, inventorySeedSQL); err != nil {
		return SafeError("inventory seed", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return SafeError("seed commit", err)
	}
	return nil
}
