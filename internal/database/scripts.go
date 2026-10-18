package database

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// SQL remains in reviewable scripts. These bindings fail at build time if a
// required script is missing; credentials are never written into SQL files.
var (
	//go:embed sql/bootstrap/01_lock.sql
	bootstrapLockSQL string
	//go:embed sql/bootstrap/02_create_role.sql
	createRoleSQL string
	//go:embed sql/bootstrap/03_configure_role.sql
	configureRoleSQL string
	//go:embed sql/bootstrap/04_role_memberships.sql
	roleMembershipsSQL string
	//go:embed sql/bootstrap/05_database_owner.sql
	databaseOwnerSQL string
	//go:embed sql/bootstrap/06_create_database.sql
	createDatabaseSQL string
	//go:embed sql/bootstrap/07_database_access.sql
	databaseAccessSQL string
	//go:embed sql/bootstrap/08_public_schema.sql
	publicSchemaSQL string
	//go:embed sql/bootstrap/09_maintenance_access.sql
	maintenanceAccessSQL string
	//go:embed sql/migrate/01_lock.sql
	migrationLockSQL string
	//go:embed sql/migrate/02_journal.sql
	migrationJournalSQL string
	//go:embed sql/migrate/03_history.sql
	migrationHistorySQL string
	//go:embed sql/migrate/04_record.sql
	migrationRecordSQL string
	//go:embed sql/migrate/05_application_access.sql
	applicationAccessSQL string
	//go:embed sql/seed/inventory.sql
	inventorySeedSQL string
)

type sqlConnection interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

// executeDDL runs the command produced by a parameterized SELECT format(...).
// PostgreSQL quotes names with %I and passwords with %L; Go never concatenates
// them into SQL. CREATE DATABASE remains outside a transaction.
func executeDDL(ctx context.Context, conn sqlConnection, query string, args ...any) error {
	var statement string
	if err := conn.QueryRow(ctx, query, args...).Scan(&statement); err != nil {
		return err
	}
	if statement == "" {
		return nil
	}
	_, err := conn.Exec(ctx, statement)
	return err
}
