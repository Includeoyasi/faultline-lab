package database

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"time"

	"github.com/Includeoyasi/faultline-lab/migrations"
	"github.com/jackc/pgx/v5"
)

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
	if _, err = tx.Exec(ctx, migrationLockSQL); err != nil {
		return SafeError("migration lock", err)
	}
	if _, err = tx.Exec(ctx, migrationJournalSQL); err != nil {
		return SafeError("migration journal", err)
	}
	rows, err := tx.Query(ctx, migrationHistorySQL)
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
		if _, err = tx.Exec(ctx, migrationRecordSQL, item.version, item.name, item.checksum); err != nil {
			return SafeError("migration recording", err)
		}
	}
	if err = executeDDL(ctx, tx, applicationAccessSQL, service.Name, service.AppUser); err != nil {
		return SafeError("application grants", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return SafeError("migration commit", err)
	}
	return nil
}
