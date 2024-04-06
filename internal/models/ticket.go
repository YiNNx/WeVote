package models

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/YiNNx/WeVote/internal/common/errors"
)

const (
	keyPrefixTicketUsage = "usage-ticket:"
)

type GlobalTicket interface {
	Access(ctx context.Context) (string, error)
	Save(ctx context.Context, ticket string) error
}

func InitGlobalTicket(ttl time.Duration) GlobalTicket {
	return &globalTicket{
		rdb: rdb,
		key: keyPrefixTicketUsage,
		ttl: ttl,
	}
}

type globalTicket struct {
	rdb redis.UniversalClient
	key string
	ttl time.Duration
}

func (t *globalTicket) Access(ctx context.Context) (string, error) {
	return t.rdb.Get(ctx, t.key).Result()
}

func (t *globalTicket) Save(ctx context.Context, ticket string) error {
	return t.rdb.Set(ctx, t.key, ticket, t.ttl).Err()
}

type Counter interface {
	Set(ctx context.Context, id string) error
	IncreaseBy(ctx context.Context, id string, count int) error
}

func InitTicketUsageCounter(ttl time.Duration, limit int) Counter {
	return &counter{
		rdb:       rdb,
		keyPrefix: "usage-ticket:",
		ttl:       ttl,
		limit:     limit,
	}
}

type counter struct {
	rdb       redis.UniversalClient
	keyPrefix string
	ttl       time.Duration
	limit     int
}

func (c *counter) key(id string) string {
	return c.keyPrefix + id
}

func (c *counter) Set(ctx context.Context, id string) error {
	return c.rdb.Set(ctx, c.key(id), 0, c.ttl).Err()
}

func (c *counter) IncreaseBy(ctx context.Context, id string, count int) error {
	n, err := c.rdb.IncrBy(ctx, c.key(id), int64(count)).Result()
	if err != nil {
		return err
	}
	if n > int64(c.limit) {
		return errors.ErrCounterLimitExceeded
	}
	return nil
}
