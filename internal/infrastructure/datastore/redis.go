package datastore

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})

	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = rdb.Ping(ctx).Result()
		cancel()

		if err == nil {
			return &RedisClient{Client: rdb}, nil
		}

		fmt.Printf(
			"failed to connect to Redis (attempt %d/%d): %v\n",
			i+1,
			maxRetries,
			err,
		)
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf(
		"failed to connect to Redis after %d attempts: %v",
		maxRetries,
		err,
	)
}

func (r *RedisClient) Set(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) Del(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *RedisClient) SAdd(
	ctx context.Context,
	key string,
	members ...interface{},
) error {
	return r.Client.SAdd(ctx, key, members...).Err()
}

func (r *RedisClient) SRem(
	ctx context.Context,
	key string,
	members ...interface{},
) error {
	return r.Client.SRem(ctx, key, members...).Err()
}

func (r *RedisClient) SMembers(
	ctx context.Context,
	key string,
) ([]string, error) {
	return r.Client.SMembers(ctx, key).Result()
}

func (r *RedisClient) Expire(
	ctx context.Context,
	key string,
	expiration time.Duration,
) error {
	return r.Client.Expire(ctx, key, expiration).Err()
}

func (r *RedisClient) Keys(
	ctx context.Context,
	pattern string,
) ([]string, error) {
	return r.Client.Keys(ctx, pattern).Result()
}
