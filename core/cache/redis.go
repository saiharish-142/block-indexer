package cache

import (
	"context"
	"time"

	"github.com/example/block-indexer/core/config"
	"github.com/redis/go-redis/v9"
)

// New returns a Redis client configured for caching.
func New(cfg config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		Password:     cfg.RedisPassword,
		DB:           0,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
	})
}

// CacheRecentTx stores tx hashes ordered by block height for an address using a sorted set.
func CacheRecentTx(ctx context.Context, rdb *redis.Client, address string, blockNumber int64, txHash string) error {
	key := "address:" + address + ":txs"
	return rdb.ZAdd(ctx, key, redis.Z{
		Score:  float64(blockNumber),
		Member: txHash,
	}).Err()
}

// FetchRecentTx retrieves the latest tx hashes for an address.
func FetchRecentTx(ctx context.Context, rdb *redis.Client, address string, limit int64) ([]string, error) {
	key := "address:" + address + ":txs"
	return rdb.ZRevRange(ctx, key, 0, limit-1).Result()
}
