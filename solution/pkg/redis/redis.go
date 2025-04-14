package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func Connect(ctx context.Context, cfg Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:   cfg.Db,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
