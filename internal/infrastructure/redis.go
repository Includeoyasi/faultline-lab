package infrastructure

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

func redisClient(cfg Config, user, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddress, Username: user, Password: password,
		Protocol: 2, DisableIdentity: true, ContextTimeoutEnabled: true,
		DialTimeout: cfg.budget(), ReadTimeout: cfg.budget(), WriteTimeout: cfg.budget(), PoolTimeout: cfg.budget(),
		MaxRetries: -1, DialerRetries: 1, PoolSize: 2, MaxActiveConns: 2,
		MaintNotificationsConfig: &maintnotifications.Config{Mode: maintnotifications.ModeDisabled},
	})
}

func redisErrorCode(err error, code string) bool {
	var serverError redis.Error
	return errors.As(err, &serverError) && strings.HasPrefix(serverError.Error(), code+" ")
}

func checkRedis(ctx context.Context, cfg Config) error {
	password, err := Password(cfg)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, cfg.budget())
	defer cancel()
	client := redisClient(cfg, cfg.RedisUser, password)
	defer func() { _ = client.Close() }()
	if err = client.Ping(ctx).Err(); err != nil {
		return safeError("Redis authenticated ping", err)
	}
	for _, auth := range []struct{ user, password, code string }{
		{"", "", "NOAUTH"}, {cfg.RedisUser, "invalid-password", "WRONGPASS"},
	} {
		denied := redisClient(cfg, auth.user, auth.password)
		err = denied.Ping(ctx).Err()
		_ = denied.Close()
		if !redisErrorCode(err, auth.code) {
			return errors.New("Redis did not return the expected authentication denial")
		}
	}
	if !redisErrorCode(client.Do(ctx, "CONFIG", "GET", "bind").Err(), "NOPERM") {
		return errors.New("Redis administrative command was not denied")
	}
	if !redisErrorCode(client.Set(ctx, "outside:"+cfg.KeyPrefix+"probe", "forbidden", 0).Err(), "NOPERM") {
		return errors.New("Redis access outside the allowed key prefix was not denied")
	}
	ttl, _ := time.ParseDuration(cfg.ProbeTTL)
	key, value := cfg.KeyPrefix+"ttl-probe", "infrastructure-check"
	if err = client.Set(ctx, key, value, ttl).Err(); err != nil {
		return safeError("Redis write", err)
	}
	got, err := client.Get(ctx, key).Result()
	if err != nil || got != value {
		return errors.New("Redis read did not match the written value")
	}
	remaining, err := client.Do(ctx, "PTTL", key).Int64()
	if err != nil || remaining <= 0 {
		return errors.New("Redis key has no positive TTL")
	}
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return safeError("Redis TTL expiration", ctx.Err())
		case <-ticker.C:
			remaining, err = client.Do(ctx, "PTTL", key).Int64()
			if err != nil {
				return safeError("Redis TTL read", err)
			}
			if remaining == -2 {
				if !errors.Is(client.Get(ctx, key).Err(), redis.Nil) {
					return errors.New("Redis expired key is still readable")
				}
				return nil
			}
		}
	}
}
