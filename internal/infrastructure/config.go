// Package infrastructure verifies the lab's Kafka and Redis dependencies.
package infrastructure

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Config is the explicit local profile; passwords are file-backed separately.
type Config struct {
	KafkaBrokers []string `json:"kafka_brokers"`
	KafkaTopic   string   `json:"kafka_topic"`
	RedisAddress string   `json:"redis_address"`
	RedisUser    string   `json:"redis_user"`
	KeyPrefix    string   `json:"key_prefix"`
	Timeout      string   `json:"timeout"`
	ProbeTTL     string   `json:"probe_ttl"`
	Secrets      string   `json:"-"`
}

var name = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,62}$`)
var topicName = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,100}$`)
var keyPrefix = regexp.MustCompile(`^[a-z][a-z0-9:._-]{0,79}:$`)

// Load reads and validates a strict JSON connection profile.
func Load(path, secrets string) (Config, error) {
	var cfg Config
	file, err := os.Open(path)
	if err != nil {
		return cfg, errors.New("cannot open infrastructure configuration")
	}
	defer func() { _ = file.Close() }()
	decoder := json.NewDecoder(io.LimitReader(file, 64<<10))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&cfg); err != nil {
		return cfg, errors.New("invalid infrastructure configuration")
	}
	var extra any
	if err = decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return cfg, errors.New("infrastructure configuration contains trailing data")
	}
	cfg.Secrets = secrets
	return cfg, cfg.Validate()
}

// Validate bounds timeouts and rejects missing addresses or unsafe ACL patterns.
func (c Config) Validate() error {
	budget, err := time.ParseDuration(c.Timeout)
	if err != nil || budget < 50*time.Millisecond || budget > 30*time.Second {
		return errors.New("infrastructure timeout must be in [50ms,30s]")
	}
	ttl, err := time.ParseDuration(c.ProbeTTL)
	if err != nil || ttl < 100*time.Millisecond || ttl >= budget {
		return errors.New("probe TTL must be at least 100ms and below the operation timeout")
	}
	if len(c.KafkaBrokers) != 1 || !address(c.KafkaBrokers[0]) || !address(c.RedisAddress) || !topicName.MatchString(c.KafkaTopic) || !name.MatchString(c.RedisUser) || !keyPrefix.MatchString(c.KeyPrefix) || c.Secrets == "" {
		return errors.New("infrastructure requires one broker, Redis address/user, topic, key prefix and secrets directory")
	}
	return nil
}

func address(value string) bool {
	host, port, err := net.SplitHostPort(value)
	if err != nil || host == "" {
		return false
	}
	number, err := strconv.ParseUint(port, 10, 16)
	return err == nil && number != 0
}

func (c Config) budget() time.Duration {
	duration, _ := time.ParseDuration(c.Timeout)
	return duration
}
