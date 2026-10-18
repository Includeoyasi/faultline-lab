package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func identifierSQL(value string) string { return pgx.Identifier{value}.Sanitize() }

// Bootstrap is repeatable after partial failure. CREATE DATABASE cannot be part
// of one transaction, so a session lock serializes this dedicated lab bootstrap.
func Bootstrap(ctx context.Context, cfg Config) error {
	admin, err := Connect(ctx, cfg, cfg.AdminDatabase, cfg.AdminUser)
	if err != nil {
		return err
	}
	defer Close(admin)
	if _, err = admin.Exec(ctx, bootstrapLockSQL); err != nil {
		return SafeError("bootstrap lock", err)
	}
	for _, service := range cfg.Services {
		for _, role := range []string{service.MigrationUser, service.AppUser} {
			if err = executeDDL(ctx, admin, createRoleSQL, role); err != nil {
				return SafeError("role creation", err)
			}
			password, err := cfg.Secret(role)
			if err != nil {
				return err
			}
			// Statement/parameter logging is disabled by the lab deployment.
			if err = executeDDL(ctx, admin, configureRoleSQL, role, password); err != nil {
				return SafeError("role configuration", err)
			}
			var memberships int
			if err = admin.QueryRow(ctx, roleMembershipsSQL, role).Scan(&memberships); err != nil {
				return SafeError("role membership check", err)
			}
			if memberships != 0 {
				return errors.New("managed roles must not inherit other roles")
			}
		}
		var owner string
		err = admin.QueryRow(ctx, databaseOwnerSQL, service.Database).Scan(&owner)
		if errors.Is(err, pgx.ErrNoRows) {
			err = executeDDL(ctx, admin, createDatabaseSQL, service.Database, service.MigrationUser)
			owner = service.MigrationUser
		}
		if err != nil {
			return SafeError("database preparation", err)
		}
		if owner != service.MigrationUser {
			return errors.New("existing database has an unexpected owner; no ownership takeover performed")
		}
		if err = executeDDL(ctx, admin, databaseAccessSQL, service.Database, service.MigrationUser, service.AppUser); err != nil {
			return SafeError("database privileges", err)
		}
		conn, err := Connect(ctx, cfg, service.Database, cfg.AdminUser)
		if err != nil {
			return err
		}
		_, err = conn.Exec(ctx, publicSchemaSQL)
		Close(conn)
		if err != nil {
			return SafeError("public schema privileges", err)
		}
	}
	for _, database := range []string{cfg.AdminDatabase, "template1"} {
		if err = executeDDL(ctx, admin, maintenanceAccessSQL, database); err != nil {
			return SafeError("maintenance database privileges", err)
		}
	}
	return nil
}
