package inventory

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config contains only Inventory's application identity. Credentials are files.
type Config struct {
	Host           string `json:"host"`
	Port           uint16 `json:"port"`
	Database       string `json:"database"`
	User           string `json:"user"`
	SSLMode        string `json:"sslmode"`
	PasswordFile   string `json:"password_file"`
	DatasetID      string `json:"-"`
	ConnectTimeout string `json:"connect_timeout"`
	RequestTimeout string `json:"request_timeout"`
	MaxConnections int32  `json:"max_connections"`
}

func LoadConfig(path, dataset string) (Config, error) {
	var cfg Config
	file, err := os.Open(path)
	if err != nil {
		return cfg, errors.New("cannot open Inventory configuration")
	}
	defer func() { _ = file.Close() }()
	decoder := json.NewDecoder(io.LimitReader(file, 16<<10))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&cfg); err != nil {
		return cfg, errors.New("invalid Inventory configuration")
	}
	var extra any
	if err = decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return cfg, errors.New("trailing Inventory configuration")
	}
	cfg.DatasetID = dataset
	return cfg, cfg.Validate()
}
func (c Config) Validate() error {
	connect, e1 := time.ParseDuration(c.ConnectTimeout)
	request, e2 := time.ParseDuration(c.RequestTimeout)
	if c.Host == "" || strings.ContainsAny(c.Host, "/\\@ \t\r\n") || c.Port == 0 || c.Database != "inventory_db" || c.User != "inventory_app" || c.SSLMode != "disable" || c.PasswordFile == "" || !identifier.MatchString(c.DatasetID) || c.MaxConnections < 1 || c.MaxConnections > 8 || e1 != nil || e2 != nil || connect < time.Millisecond || connect > time.Minute || request < time.Millisecond || request > time.Minute {
		return errors.New("inventory requires its application identity, dataset, secret file, explicit endpoint, bounded timeouts and 1..8 connections")
	}
	return nil
}
func Open(ctx context.Context, c Config) (*Store, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	password, err := os.ReadFile(c.PasswordFile)
	if err != nil {
		return nil, errors.New("inventory credential unavailable")
	}
	secret := strings.TrimSpace(string(password))
	if len(secret) != 64 {
		return nil, errors.New("inventory credential malformed")
	}
	for _, char := range secret {
		if !strings.ContainsRune("0123456789abcdef", char) {
			return nil, errors.New("inventory credential malformed")
		}
	}
	address := url.URL{Scheme: "postgres", Host: net.JoinHostPort(c.Host, strconv.Itoa(int(c.Port))), Path: c.Database, User: url.UserPassword(c.User, secret)}
	address.RawQuery = url.Values{"sslmode": {c.SSLMode}, "require_auth": {"scram-sha-256"}}.Encode()
	cfg, err := pgxpool.ParseConfig(address.String())
	if err != nil {
		return nil, errors.New("cannot configure Inventory database")
	}
	cfg.MaxConns = c.MaxConnections
	cfg.MinConns = 0
	cfg.ConnConfig.ConnectTimeout, _ = time.ParseDuration(c.ConnectTimeout)
	budget, _ := time.ParseDuration(c.RequestTimeout)
	cfg.ConnConfig.RuntimeParams["statement_timeout"] = strconv.FormatInt(budget.Milliseconds(), 10)
	cfg.ConnConfig.RuntimeParams["lock_timeout"] = strconv.FormatInt(budget.Milliseconds(), 10)
	cfg.ConnConfig.RuntimeParams["timezone"] = "UTC"
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, storageError(err)
	}
	store, err := NewStore(ctx, pool, c.DatasetID, budget)
	if err != nil {
		pool.Close()
		return nil, err
	}
	return store, nil
}
