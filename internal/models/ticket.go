package models

import (
	"context"
	"time"

	"github.com/YiNNx/WeVote/internal/common/errors"
	"github.com/redis/go-redis/v9"
)

type Counter struct {
	rdb       redis.UniversalClient
	keyPrefix string
	ttl       time.Duration
	limit     int
}

func (c *Counter) key(id string) string {
	return c.keyPrefix + id
}

func (c *Counter) Set(ctx context.Context, id string) error {
	return c.rdb.Set(ctx, c.key(id), 0, c.ttl).Err()
}

func (c *Counter) IncreaseBy(ctx context.Context, id string, count int) error {
	n, err := c.rdb.IncrBy(ctx, c.key(id), int64(count)).Result()
	if err != nil {
		return err
	}
	if n > int64(c.limit) {
		return errors.ErrCounterLimitExceeded
	}
	return nil
}
