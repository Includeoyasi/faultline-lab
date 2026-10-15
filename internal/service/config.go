// Package service contains the shared startup mechanics of the three services.
package service

import (
	"errors"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"
)

// Config holds process settings. ShutdownTimeout is reserved for P01-05.
type Config struct {
	GRPCAddress     string
	LogLevel        slog.Level
	ShutdownTimeout time.Duration
}

// LoadConfig applies defaults and validates environment overrides without
// including the supplied values in errors. An explicitly empty value is invalid.
func LoadConfig(defaultAddress string, lookup func(string) (string, bool)) (Config, error) {
	cfg := Config{GRPCAddress: defaultAddress, LogLevel: slog.LevelInfo, ShutdownTimeout: 10 * time.Second}
	if value, ok := lookup("FL_GRPC_ADDR"); ok {
		cfg.GRPCAddress = value
	}
	if !validAddress(cfg.GRPCAddress) {
		return Config{}, errors.New("FL_GRPC_ADDR must be host:port with a port from 1 to 65535")
	}
	if value, ok := lookup("FL_LOG_LEVEL"); ok {
		switch strings.ToUpper(value) {
		case "DEBUG":
			cfg.LogLevel = slog.LevelDebug
		case "INFO":
			cfg.LogLevel = slog.LevelInfo
		case "WARN":
			cfg.LogLevel = slog.LevelWarn
		case "ERROR":
			cfg.LogLevel = slog.LevelError
		default:
			return Config{}, errors.New("FL_LOG_LEVEL must be DEBUG, INFO, WARN or ERROR")
		}
	}
	if value, ok := lookup("FL_SHUTDOWN_TIMEOUT"); ok {
		timeout, err := time.ParseDuration(value)
		if err != nil || timeout <= 0 {
			return Config{}, errors.New("FL_SHUTDOWN_TIMEOUT must be a positive Go duration")
		}
		cfg.ShutdownTimeout = timeout
	}
	return cfg, nil
}

func validAddress(address string) bool {
	host, port, err := net.SplitHostPort(address)
	if err != nil || port == "" {
		return false
	}
	for _, char := range port {
		if char < '0' || char > '9' {
			return false
		}
	}
	number, err := strconv.Atoi(port)
	if err != nil || number < 1 || number > 65535 {
		return false
	}
	if host == "" || net.ParseIP(host) != nil {
		return true
	}
	// Validate DNS syntax locally; startup must not require a DNS lookup.
	host = strings.TrimSuffix(host, ".")
	if len(host) == 0 || len(host) > 253 {
		return false
	}
	for _, label := range strings.Split(host, ".") {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for _, char := range label {
			switch {
			case char >= 'a' && char <= 'z', char >= 'A' && char <= 'Z', char >= '0' && char <= '9', char == '-':
			default:
				return false
			}
		}
	}
	return true
}
