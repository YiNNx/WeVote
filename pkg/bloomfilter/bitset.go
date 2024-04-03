package bloom

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisBitSet struct {
	rdb redis.UniversalClient
	key string
	m   uint
}

func NewRedisBitSet(rdb redis.UniversalClient, key string, m uint) *RedisBitSet {
	return &RedisBitSet{rdb, key, m}
}

func (b *RedisBitSet) Set(ctx context.Context, offsets []uint) error {
	rtx := b.rdb.Pipeline()
	for _, offset := range offsets {
		_, err := rtx.SetBit(ctx, b.key, int64(offset), 1).Result()
		if err != nil {
			rtx.Discard()
			return err
		}
	}
	_, err := rtx.Exec(ctx)
	return err
}

func (b *RedisBitSet) Test(ctx context.Context, offsets []uint) (bool, error) {
	for _, offset := range offsets {
		res, err := b.rdb.GetBit(ctx, b.key, int64(offset)).Result()
		if err != nil {
			return false, err
		}
		if res == 0 {
			return false, nil
		}
	}

	return true, nil
}
