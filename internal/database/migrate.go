package database

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"strconv"
	"time"

	"github.com/Includeoyasi/faultline-lab/migrations"
	"github.com/jackc/pgx/v5"
)

type migration struct {
	version             int
	name, checksum, sql string
}

var migrationName = regexp.MustCompile(`^([0-9]{4})_[a-z0-9_]+\.sql$`)

func loadMigrations(files fs.FS, service string) ([]migration, error) {
	entries, err := fs.ReadDir(files, service)
	if err != nil {
		return nil, errors.New("cannot read migrations")
	}
	var result []migration
	previous := 0
	for _, entry := range entries {
		match := migrationName.FindStringSubmatch(entry.Name())
		if entry.IsDir() || match == nil {
			return nil, errors.New("invalid migration filename")
		}
		version, _ := strconv.Atoi(match[1])
		if version <= previous {
			return nil, errors.New("migration versions must be positive and strictly increasing")
		}
		data, err := fs.ReadFile(files, service+"/"+entry.Name())
		if err != nil || len(data) == 0 {
			return nil, errors.New("cannot read nonempty migration")
		}
		hash := sha256.Sum256(data)
		result = append(result, migration{version, entry.Name(), hex.EncodeToString(hash[:]), string(data)})
		previous = version
	}
	if len(result) == 0 {
		return nil, errors.New("service has no migrations")
	}
	return result, nil
}

// Migrate atomically applies pending SQL and verifies immutable service history.
func Migrate(ctx context.Context, cfg Config, service Service) error {
	return migrate(ctx, cfg, service, migrations.Files)
}

func rollback(tx pgx.Tx) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = tx.Rollback(ctx)
}

func migrate(ctx context.Context, cfg Config, service Service, files fs.FS) error {
	items, err := loadMigrations(files, service.Name)
	if err != nil {
		return err
	}
	conn, err := Connect(ctx, cfg, service.Database, service.MigrationUser)
	if err != nil {
		return err
	}
	defer Close(conn)
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return SafeError("migration transaction", err)
	}
	defer rollback(tx)
	// Lock before creating or reading the journal. It is local to this database.
	if _, err = tx.Exec(ctx, `SELECT pg_advisory_xact_lock(706020121)`); err != nil {
		return SafeError("migration lock", err)
	}
	if _, err = tx.Exec(ctx, `CREATE SCHEMA IF NOT EXISTS lab_migrations;
CREATE TABLE IF NOT EXISTS lab_migrations.versions (
 version integer PRIMARY KEY, name text NOT NULL, checksum text NOT NULL,
 applied_at timestamptz NOT NULL DEFAULT now()
)`); err != nil {
		return SafeError("migration journal", err)
	}
	rows, err := tx.Query(ctx, "SELECT version,name,checksum FROM lab_migrations.versions ORDER BY version")
	if err != nil {
		return SafeError("migration history", err)
	}
	var applied []migration
	for rows.Next() {
		var item migration
		if err = rows.Scan(&item.version, &item.name, &item.checksum); err != nil {
			rows.Close()
			return SafeError("migration history", err)
		}
		applied = append(applied, item)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return SafeError("migration history", err)
	}
	if len(applied) > len(items) {
		return errors.New("database migration history is ahead of this tool")
	}
	for index, item := range items {
		if index < len(applied) {
			if old := applied[index]; old.version != item.version || old.name != item.name || old.checksum != item.checksum {
				return fmt.Errorf("migration history/checksum mismatch at version %d", item.version)
			}
			continue
		}
		if _, err = tx.Exec(ctx, item.sql); err != nil {
			return SafeError(fmt.Sprintf("migration %d", item.version), err)
		}
		if _, err = tx.Exec(ctx, "INSERT INTO lab_migrations.versions(version,name,checksum) VALUES($1,$2,$3)", item.version, item.name, item.checksum); err != nil {
			return SafeError("migration recording", err)
		}
	}
	schema, app := identifierSQL(service.Name), identifierSQL(service.AppUser)
	for _, statement := range []string{
		"GRANT USAGE ON SCHEMA " + schema + " TO " + app,
		"GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA " + schema + " TO " + app,
		"ALTER DEFAULT PRIVILEGES IN SCHEMA " + schema + " GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO " + app,
	} {
		if _, err = tx.Exec(ctx, statement); err != nil {
			return SafeError("application grants", err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return SafeError("migration commit", err)
	}
	return nil
}

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
	data, err := migrations.Files.ReadFile("seed/inventory.sql")
	if err != nil {
		return errors.New("cannot read inventory seed")
	}
	if _, err = tx.Exec(ctx, string(data)); err != nil {
		return SafeError("inventory seed", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return SafeError("seed commit", err)
	}
	return nil
}
