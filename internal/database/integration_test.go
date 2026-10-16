package database

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	"github.com/Includeoyasi/faultline-lab/migrations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// These tests change their dedicated dataset. Only the isolated host smoke opts in.
func TestPostgresFoundation(t *testing.T) {
	if os.Getenv("FL_DATA_INTEGRATION") != "1" {
		t.Skip("requires isolated PostgreSQL smoke")
	}
	cfg, err := Load(os.Getenv("FL_DATA_CONFIG"), os.Getenv("FL_DATA_SECRETS"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	connect := func(t *testing.T, db, user string) *pgx.Conn {
		t.Helper()
		conn, e := Connect(ctx, cfg, db, user)
		if e != nil {
			t.Fatal(e)
		}
		t.Cleanup(func() { Close(conn) })
		return conn
	}
	scalar := func(t *testing.T, conn *pgx.Conn, query string) int {
		t.Helper()
		var value int
		if e := conn.QueryRow(ctx, query).Scan(&value); e != nil {
			t.Fatal(SafeError("test query", e))
		}
		return value
	}
	denied := func(t *testing.T, conn *pgx.Conn, sql string) {
		t.Helper()
		_, e := conn.Exec(ctx, sql)
		var pgerr *pgconn.PgError
		if !errors.As(e, &pgerr) || pgerr.Code != "42501" {
			t.Fatalf("expected permission denial, got %v", SafeError("negative permission test", e))
		}
	}
	t.Run("concurrent migrations and immutable history", func(t *testing.T) {
		for _, service := range cfg.Services {
			var wg sync.WaitGroup
			results := make(chan error, 2)
			for attempt := 0; attempt < 2; attempt++ {
				wg.Go(func() { results <- Migrate(ctx, cfg, service) })
			}
			wg.Wait()
			close(results)
			for e := range results {
				if e != nil {
					t.Fatal(e)
				}
			}
			conn := connect(t, service.Database, service.MigrationUser)
			want := 1
			if service.Name == "inventory" {
				want = 2
			}
			if got := scalar(t, conn, "SELECT count(*) FROM lab_migrations.versions"); got != want {
				t.Fatalf("versions=%d want=%d", got, want)
			}
		}
		before, e := Fingerprint(ctx, cfg)
		if e != nil {
			t.Fatal(e)
		}
		for _, service := range cfg.Services {
			if e = Migrate(ctx, cfg, service); e != nil {
				t.Fatal(e)
			}
		}
		after, e := Fingerprint(ctx, cfg)
		if e != nil || before != after {
			t.Fatal("repeat modified history", e)
		}
	})
	t.Run("role isolation and denied DDL", func(t *testing.T) {
		for _, service := range cfg.Services {
			for _, user := range []string{service.AppUser, service.MigrationUser} {
				for _, target := range cfg.Services {
					conn, e := Connect(ctx, cfg, target.Database, user)
					if target.Name == service.Name {
						if e != nil {
							t.Fatal(e)
						}
						Close(conn)
					} else {
						if conn != nil {
							Close(conn)
						}
						if e == nil || !strings.Contains(e.Error(), "42501") {
							t.Fatalf("cross-database access allowed or wrong error: %v", e)
						}
					}
				}
			}
			app := connect(t, service.Database, service.AppUser)
			for _, sql := range []string{"CREATE SCHEMA forbidden", "CREATE TABLE public.forbidden(id int)", "CREATE TEMP TABLE forbidden(id int)", "SET ROLE " + identifierSQL(service.MigrationUser), "DELETE FROM lab_migrations.versions"} {
				denied(t, app, sql)
			}
			if got := scalar(t, app, `SELECT count(*) FROM pg_roles WHERE rolname=current_user AND (rolsuper OR rolcreatedb OR rolcreaterole OR rolreplication OR rolbypassrls)`); got != 0 {
				t.Fatal("elevated app role")
			}
			migration := connect(t, service.Database, service.MigrationUser)
			table := identifierSQL(service.Name) + ".privilege_probe"
			if _, e := migration.Exec(ctx, "CREATE TABLE "+table+"(id int PRIMARY KEY)"); e != nil {
				t.Fatal(SafeError("probe create", e))
			}
			for _, sql := range []string{"INSERT INTO " + table + " VALUES(1)", "UPDATE " + table + " SET id=2", "DELETE FROM " + table + " WHERE id=2"} {
				if _, e := app.Exec(ctx, sql); e != nil {
					t.Fatal(SafeError("app DML", e))
				}
			}
			denied(t, app, "ALTER TABLE "+table+" ADD COLUMN forbidden int")
			denied(t, app, "TRUNCATE "+table)
			if _, e := migration.Exec(ctx, "DROP TABLE "+table); e != nil {
				t.Fatal(SafeError("probe drop", e))
			}
			for _, maintenance := range []string{cfg.AdminDatabase, "template1"} {
				conn, e := Connect(ctx, cfg, maintenance, service.AppUser)
				if conn != nil {
					Close(conn)
				}
				if e == nil || !strings.Contains(e.Error(), "42501") {
					t.Fatalf("maintenance access: %v", e)
				}
			}
		}
	})
	t.Run("failed migration rollback and checksum drift", func(t *testing.T) {
		service, _ := cfg.Service("inventory")
		conn := connect(t, service.Database, service.MigrationUser)
		files := fstest.MapFS{}
		entries, e := fs.ReadDir(migrations.Files, "inventory")
		if e != nil {
			t.Fatal(e)
		}
		for _, entry := range entries {
			data, e := fs.ReadFile(migrations.Files, "inventory/"+entry.Name())
			if e != nil {
				t.Fatal(e)
			}
			files["inventory/"+entry.Name()] = &fstest.MapFile{Data: data}
		}
		files["inventory/0003_failure.sql"] = &fstest.MapFile{Data: []byte("CREATE TABLE inventory.must_rollback(id int); SELECT 1/0;")}
		if e = migrate(ctx, cfg, service, files); e == nil || !strings.Contains(e.Error(), "22012") {
			t.Fatalf("expected failed migration: %v", e)
		}
		if scalar(t, conn, "SELECT count(*) FROM lab_migrations.versions") != 2 || scalar(t, conn, "SELECT count(*) FROM information_schema.tables WHERE table_schema='inventory' AND table_name='must_rollback'") != 0 {
			t.Fatal("failed migration persisted")
		}
		delete(files, "inventory/0003_failure.sql")
		files["inventory/0001_schema.sql"].Data = []byte("CREATE SCHEMA inventory; -- changed")
		if e = migrate(ctx, cfg, service, files); e == nil || !strings.Contains(e.Error(), "checksum") {
			t.Fatal("accepted changed migration", e)
		}
	})
	t.Run("repeat seed preserves catalog and modified stock", func(t *testing.T) {
		if e := Seed(ctx, cfg); e != nil {
			t.Fatal(e)
		}
		service, _ := cfg.Service("inventory")
		conn := connect(t, service.Database, service.AppUser)
		if scalar(t, conn, "SELECT count(*) FROM inventory.products") != 3 || scalar(t, conn, "SELECT available FROM inventory.stock WHERE sku='SKU-001'") != 100 {
			t.Fatal("unexpected initial catalog")
		}
		if _, e := conn.Exec(ctx, "UPDATE inventory.stock SET available=7 WHERE sku='SKU-001'; UPDATE inventory.products SET price_minor=99900 WHERE sku='SKU-001'"); e != nil {
			t.Fatal(SafeError("test update", e))
		}
		if e := Seed(ctx, cfg); e != nil {
			t.Fatal(e)
		}
		if scalar(t, conn, "SELECT count(*) FROM inventory.products") != 3 || scalar(t, conn, "SELECT available FROM inventory.stock WHERE sku='SKU-001'") != 7 || scalar(t, conn, "SELECT price_minor FROM inventory.products WHERE sku='SKU-001'") != 99900 {
			t.Fatal("repeat seed overwrote data")
		}
		for _, sql := range []string{"UPDATE inventory.stock SET available=-1 WHERE sku='SKU-001'", "UPDATE inventory.products SET price_minor=-1 WHERE sku='SKU-001'", "UPDATE inventory.products SET currency='rub' WHERE sku='SKU-001'", "INSERT INTO inventory.stock VALUES('missing',1)"} {
			if _, e := conn.Exec(ctx, sql); e == nil {
				t.Fatal("invalid inventory data accepted")
			}
		}
	})
}
