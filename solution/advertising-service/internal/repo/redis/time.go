package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type TimeRepo struct {
	rdb *redis.Client
}

func NewTimeRepo(rdb *redis.Client) *TimeRepo {
	return &TimeRepo{
		rdb: rdb,
	}
}

func (tr *TimeRepo) SetDay(ctx context.Context, date int) error {
	op := "TimeRepo.SetDay"

	err := tr.rdb.Set(ctx, tr.getKey(), date, -1).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (tr *TimeRepo) GetDay(ctx context.Context) (int, error) {
	op := "TimeRepo.GetDay"

	date, err := tr.rdb.Get(ctx, tr.getKey()).Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return date, nil
}

func (tr *TimeRepo) getKey() string {
	return "time.current-day"
}
