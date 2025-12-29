package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/example/block-indexer/pkg/config"
)

// NewRedisClient builds a redis client from config.
func NewRedisClient(cfg config.RedisConfig) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	}
	client := redis.NewClient(opts)
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("connect redis: %w", err)
	}
	return client, nil
}
