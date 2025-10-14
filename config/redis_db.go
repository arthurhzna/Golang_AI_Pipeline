package config

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func CreateClient(ctx context.Context, dbNo int, addr string, password string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbNo,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis at %s: %w", addr, err)
	}

	return rdb, nil
}
