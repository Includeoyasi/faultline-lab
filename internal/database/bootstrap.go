package database

import (
	"context"
	"errors"
	"strings"

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
	if _, err = admin.Exec(ctx, "SELECT pg_advisory_lock(706020120)"); err != nil {
		return SafeError("bootstrap lock", err)
	}
	for _, service := range cfg.Services {
		for _, role := range []string{service.MigrationUser, service.AppUser} {
			var exists bool
			if err = admin.QueryRow(ctx, "SELECT EXISTS(SELECT FROM pg_roles WHERE rolname=$1)", role).Scan(&exists); err != nil {
				return SafeError("role lookup", err)
			}
			if !exists {
				if _, err = admin.Exec(ctx, "CREATE ROLE "+identifierSQL(role)); err != nil {
					return SafeError("role creation", err)
				}
			}
			password, err := cfg.Secret(role)
			if err != nil {
				return err
			}
			// Passwords are generated hex, and still quoted defensively. Server
			// statement/parameter logging is disabled by the lab deployment.
			statement := "ALTER ROLE " + identifierSQL(role) + " WITH LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS NOINHERIT PASSWORD '" + strings.ReplaceAll(password, "'", "''") + "'"
			if _, err = admin.Exec(ctx, statement); err != nil {
				return SafeError("role configuration", err)
			}
			var memberships int
			if err = admin.QueryRow(ctx, "SELECT count(*) FROM pg_auth_members m JOIN pg_roles r ON r.oid=m.member WHERE r.rolname=$1", role).Scan(&memberships); err != nil {
				return SafeError("role membership check", err)
			}
			if memberships != 0 {
				return errors.New("managed roles must not inherit other roles")
			}
		}
		var owner string
		err = admin.QueryRow(ctx, "SELECT pg_get_userbyid(datdba) FROM pg_database WHERE datname=$1", service.Database).Scan(&owner)
		if errors.Is(err, pgx.ErrNoRows) {
			_, err = admin.Exec(ctx, "CREATE DATABASE "+identifierSQL(service.Database)+" OWNER "+identifierSQL(service.MigrationUser))
			owner = service.MigrationUser
		}
		if err != nil {
			return SafeError("database preparation", err)
		}
		if owner != service.MigrationUser {
			return errors.New("existing database has an unexpected owner; no ownership takeover performed")
		}
		for _, statement := range []string{
			"REVOKE ALL ON DATABASE " + identifierSQL(service.Database) + " FROM PUBLIC",
			"GRANT CONNECT ON DATABASE " + identifierSQL(service.Database) + " TO " + identifierSQL(service.MigrationUser) + ", " + identifierSQL(service.AppUser),
		} {
			if _, err = admin.Exec(ctx, statement); err != nil {
				return SafeError("database privileges", err)
			}
		}
		conn, err := Connect(ctx, cfg, service.Database, cfg.AdminUser)
		if err != nil {
			return err
		}
		_, err = conn.Exec(ctx, "REVOKE ALL ON SCHEMA public FROM PUBLIC")
		Close(conn)
		if err != nil {
			return SafeError("public schema privileges", err)
		}
	}
	for _, database := range []string{cfg.AdminDatabase, "template1"} {
		if _, err = admin.Exec(ctx, "REVOKE ALL ON DATABASE "+identifierSQL(database)+" FROM PUBLIC"); err != nil {
			return SafeError("maintenance database privileges", err)
		}
	}
	return nil
}
