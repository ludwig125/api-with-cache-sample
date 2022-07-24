package main

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	r      *redis.Client
	maxTTL time.Duration
}

func NewRedisCache(maxTTL int) CacheRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisCache{
		r:      rdb,
		maxTTL: time.Duration(maxTTL * int(time.Second)),
	}
}

func (c *RedisCache) Put(ctx context.Context, k ID, v bool) error {
	return c.r.Set(ctx, k.String(), v, c.maxTTL).Err()
}

func (c *RedisCache) Get(ctx context.Context, k ID) (bool, error) {
	// val, err := rdb.Get(ctx, "key").Result()
	_, err := c.r.Get(ctx, k.String()).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
