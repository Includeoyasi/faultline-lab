// Package database prepares the lab's PostgreSQL databases without business RPCs.
package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Service assigns database and role ownership to one business service.
type Service struct {
	Name          string `json:"name"`
	Database      string `json:"database"`
	AppUser       string `json:"app_user"`
	MigrationUser string `json:"migration_user"`
}

// Config contains the explicit local connection profile; credentials stay in files.
type Config struct {
	Host          string    `json:"host"`
	Port          uint16    `json:"port"`
	SSLMode       string    `json:"sslmode"`
	Timeout       string    `json:"timeout"`
	AdminDatabase string    `json:"admin_database"`
	AdminUser     string    `json:"admin_user"`
	Services      []Service `json:"services"`
	Secrets       string    `json:"-"`
}

var identifier = regexp.MustCompile(`^[a-z][a-z0-9_]{0,62}$`)
var secretValue = regexp.MustCompile(`^[a-f0-9]{64}$`)

// Load reads a strict JSON profile and validates it before any connection.
func Load(path, secrets string) (Config, error) {
	var cfg Config
	file, err := os.Open(path)
	if err != nil {
		return cfg, errors.New("cannot open database configuration")
	}
	defer func() { _ = file.Close() }()
	decoder := json.NewDecoder(io.LimitReader(file, 64<<10))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&cfg); err != nil {
		return cfg, errors.New("invalid database configuration")
	}
	var extra any
	if err = decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return cfg, errors.New("database configuration contains trailing data")
	}
	cfg.Secrets = secrets
	return cfg, cfg.Validate()
}

// Validate rejects incomplete profiles, shared identities and unbounded timeouts.
func (c Config) Validate() error {
	budget, err := time.ParseDuration(c.Timeout)
	if err != nil || budget < time.Millisecond || budget > time.Minute || c.Host == "" || c.Port == 0 || c.SSLMode != "disable" || c.Secrets == "" {
		return errors.New("database configuration requires host, port, local sslmode=disable, secrets directory and timeout in [1ms,1m]")
	}
	if !identifier.MatchString(c.AdminDatabase) || !identifier.MatchString(c.AdminUser) || len(c.Services) != 3 {
		return errors.New("invalid administrator or service configuration")
	}
	names, databases, users := map[string]bool{}, map[string]bool{c.AdminDatabase: true}, map[string]bool{c.AdminUser: true}
	for _, service := range c.Services {
		if (service.Name != "orders" && service.Name != "inventory" && service.Name != "payments") || names[service.Name] || !identifier.MatchString(service.Database) || databases[service.Database] {
			return errors.New("configuration must contain three distinct service databases")
		}
		names[service.Name], databases[service.Database] = true, true
		for _, user := range []string{service.AppUser, service.MigrationUser} {
			if !identifier.MatchString(user) || users[user] {
				return errors.New("database roles must be valid and distinct")
			}
			users[user] = true
		}
	}
	return nil
}

// Service finds one of the three configured service identities.
func (c Config) Service(name string) (Service, error) {
	for _, service := range c.Services {
		if service.Name == name {
			return service, nil
		}
	}
	return Service{}, errors.New("unknown service")
}

// Secret reads a generated hex password without including its value in errors.
func (c Config) Secret(user string) (string, error) {
	if !identifier.MatchString(user) {
		return "", errors.New("invalid secret identifier")
	}
	data, err := os.ReadFile(filepath.Join(c.Secrets, user+".password"))
	if err != nil || !secretValue.MatchString(strings.TrimSpace(string(data))) {
		return "", errors.New("database secret is unavailable or malformed")
	}
	return strings.TrimSpace(string(data)), nil
}

// GenerateSecrets preserves existing credentials. The private directory protects
// host access; readable files allow non-root containers to consume bind secrets.
func GenerateSecrets(cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	if err := os.MkdirAll(cfg.Secrets, 0o700); err != nil {
		return errors.New("cannot create private secret directory")
	}
	info, err := os.Lstat(cfg.Secrets)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return errors.New("secret directory must be a regular directory")
	}
	if err = os.Chmod(cfg.Secrets, 0o700); err != nil {
		return errors.New("cannot protect secret directory")
	}
	users := []string{cfg.AdminUser}
	for _, service := range cfg.Services {
		users = append(users, service.AppUser, service.MigrationUser)
	}
	for _, user := range users {
		path := filepath.Join(cfg.Secrets, user+".password")
		if info, err := os.Lstat(path); err == nil {
			if !info.Mode().IsRegular() {
				return errors.New("secret must be a regular file")
			}
			if _, err = cfg.Secret(user); err != nil {
				return err
			}
			if err = os.Chmod(path, 0o644); err != nil {
				return errors.New("cannot set secret file permissions")
			}
			continue
		} else if !errors.Is(err, os.ErrNotExist) {
			return errors.New("cannot inspect secret file")
		}
		var value [32]byte
		if _, err := rand.Read(value[:]); err != nil {
			return errors.New("cannot generate database secret")
		}
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			return errors.New("cannot create secret file")
		}
		_, writeErr := file.WriteString(hex.EncodeToString(value[:]) + "\n")
		permissionErr := file.Chmod(0o644) // Creation mode alone is restricted by POSIX umask.
		if closeErr := file.Close(); errors.Join(writeErr, permissionErr, closeErr) != nil {
			return errors.New("cannot save secret file")
		}
	}
	return nil
}

// Connect uses explicit credentials, SCRAM authentication and bounded SQL waits.
func Connect(ctx context.Context, cfg Config, database, user string) (*pgx.Conn, error) {
	password, err := cfg.Secret(user)
	if err != nil {
		return nil, err
	}
	address := url.URL{Scheme: "postgres", Host: net.JoinHostPort(cfg.Host, strconv.Itoa(int(cfg.Port))), Path: database, User: url.UserPassword(user, password)}
	query := url.Values{"sslmode": {cfg.SSLMode}, "require_auth": {"scram-sha-256"}}
	address.RawQuery = query.Encode()
	connection, err := pgx.ParseConfig(address.String())
	if err != nil {
		return nil, errors.New("cannot configure database connection")
	}
	connection.ConnectTimeout, _ = time.ParseDuration(cfg.Timeout)
	connection.RuntimeParams["statement_timeout"] = strconv.FormatInt(connection.ConnectTimeout.Milliseconds(), 10)
	connection.RuntimeParams["lock_timeout"] = strconv.FormatInt(connection.ConnectTimeout.Milliseconds(), 10)
	connection.RuntimeParams["timezone"] = "UTC"
	conn, err := pgx.ConnectConfig(ctx, connection)
	if err != nil {
		return nil, SafeError("database connection", err)
	}
	return conn, nil
}

// Close bounds connection cleanup independently of the operation's context.
func Close(conn *pgx.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = conn.Close(ctx)
}

// SafeError excludes SQL, parameters, connection strings and server details.
func SafeError(stage string, err error) error {
	var databaseError *pgconn.PgError
	if errors.As(err, &databaseError) {
		return fmt.Errorf("%s failed (SQLSTATE %s)", stage, databaseError.Code)
	}
	if errors.Is(err, context.Canceled) {
		return fmt.Errorf("%s: %w", stage, context.Canceled)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%s: %w", stage, context.DeadlineExceeded)
	}
	return fmt.Errorf("%s failed", stage)
}
