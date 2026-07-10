package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/config"
)

// Client defines the cache operations used by application services.
type Client interface {
	Get(ctx context.Context, key string, dest interface{}) (bool, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	DeleteByPrefix(ctx context.Context, prefix string) error
	IsEnabled() bool
	Close() error
}

type redisClient struct {
	enabled bool
	client  *redis.Client
	log     *zap.Logger
}

// NewClient creates a Redis-backed cache client when enabled,
// otherwise it returns a disabled no-op implementation.
func NewClient(cfg *config.RedisConfig, log *zap.Logger) (Client, error) {
	if cfg == nil || !cfg.Enabled {
		return &redisClient{enabled: false, log: log}, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Info("redis cache connected",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.Int("db", cfg.DB),
	)

	return &redisClient{
		enabled: true,
		client:  client,
		log:     log,
	}, nil
}

func (c *redisClient) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	if !c.enabled {
		return false, nil
	}

	raw, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get cache key %q: %w", key, err)
	}

	if err := json.Unmarshal([]byte(raw), dest); err != nil {
		return false, fmt.Errorf("failed to decode cache key %q: %w", key, err)
	}

	return true, nil
}

func (c *redisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if !c.enabled {
		return nil
	}

	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to encode cache key %q: %w", key, err)
	}

	if err := c.client.Set(ctx, key, payload, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache key %q: %w", key, err)
	}

	return nil
}

func (c *redisClient) Delete(ctx context.Context, keys ...string) error {
	if !c.enabled || len(keys) == 0 {
		return nil
	}

	if err := c.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("failed to delete cache keys: %w", err)
	}

	return nil
}

func (c *redisClient) DeleteByPrefix(ctx context.Context, prefix string) error {
	if !c.enabled {
		return nil
	}

	pattern := prefix + "*"
	var cursor uint64

	for {
		keys, nextCursor, err := c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("failed to scan cache keys with prefix %q: %w", prefix, err)
		}

		if len(keys) > 0 {
			if err := c.client.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("failed to delete cache keys with prefix %q: %w", prefix, err)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}

func (c *redisClient) IsEnabled() bool {
	return c.enabled
}

func (c *redisClient) Close() error {
	if !c.enabled || c.client == nil {
		return nil
	}
	return c.client.Close()
}
