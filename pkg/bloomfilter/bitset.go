package bloom

import (
	"context"
	"math"

	"github.com/redis/go-redis/v9"
)

var redisBitSetMaxLength = uint(math.Pow(2, 32))

type RedisBitSet struct {
	rdb redis.UniversalClient
	key string
}

func NewRedisBitSet(rdb redis.UniversalClient, key string) *RedisBitSet {
	return &RedisBitSet{rdb, key}
}

func (b *RedisBitSet) Set(ctx context.Context, offsets []uint) error {
	rtx := b.rdb.Pipeline()
	for _, offset := range offsets {
		_, err := rtx.SetBit(ctx, b.key, int64(offset/redisBitSetMaxLength), 1).Result()
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
		res, err := b.rdb.GetBit(ctx, b.key, int64(offset/redisBitSetMaxLength)).Result()
		if err != nil || res == 0 {
			return false, err
		}
	}
	return true, nil
}
