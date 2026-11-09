package store

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

	"github.com/Includeoyasi/faultline-lab/internal/orders/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Config contains only Orders's application identity. Credentials are files.
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
		return cfg, errors.New("cannot open Orders configuration")
	}
	defer func() { _ = file.Close() }()
	decoder := json.NewDecoder(io.LimitReader(file, 16<<10))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&cfg); err != nil {
		return cfg, errors.New("invalid Orders configuration")
	}
	var extra any
	if err = decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return cfg, errors.New("trailing Orders configuration")
	}
	cfg.DatasetID = dataset
	return cfg, cfg.Validate()
}
func (c Config) Validate() error {
	connect, e1 := time.ParseDuration(c.ConnectTimeout)
	request, e2 := time.ParseDuration(c.RequestTimeout)
	if c.Host == "" || strings.ContainsAny(c.Host, "/\\@ \t\r\n") || c.Port == 0 || c.Database != "orders_db" || c.User != "orders_app" || c.SSLMode != "disable" || c.PasswordFile == "" || !model.ValidIdentifier(c.DatasetID) || c.MaxConnections < 1 || c.MaxConnections > 8 || e1 != nil || e2 != nil || connect < time.Millisecond || connect > time.Minute || request < time.Millisecond || request > time.Minute {
		return errors.New("orders requires its application identity, dataset, secret file, explicit endpoint, bounded timeouts and 1..8 connections")
	}
	return nil
}
func Open(ctx context.Context, c Config) (*Store, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	password, err := os.ReadFile(c.PasswordFile)
	if err != nil {
		return nil, errors.New("orders credential unavailable")
	}
	secret := strings.TrimSpace(string(password))
	if len(secret) != 64 {
		return nil, errors.New("orders credential malformed")
	}
	for _, char := range secret {
		if !strings.ContainsRune("0123456789abcdef", char) {
			return nil, errors.New("orders credential malformed")
		}
	}
	address := url.URL{Scheme: "postgres", Host: net.JoinHostPort(c.Host, strconv.Itoa(int(c.Port))), Path: c.Database, User: url.UserPassword(c.User, secret)}
	address.RawQuery = url.Values{"sslmode": {c.SSLMode}, "require_auth": {"scram-sha-256"}}.Encode()
	cfg, err := pgxpool.ParseConfig(address.String())
	if err != nil {
		return nil, errors.New("cannot configure Orders database")
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
