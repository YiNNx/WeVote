package models

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/YiNNx/WeVote/pkg/bloomfilter"
)

const (
	redisBitSetMaxLength = 4 * 1024 * 1024 * 1024
)

type redisBitSet struct {
	rdb redis.UniversalClient
	key string
}

func NewRedisBitSet(key string) bloomfilter.BitSetProvider {
	return &redisBitSet{rdb, key}
}

func (b *redisBitSet) Set(ctx context.Context, offsets []uint) error {
	rtx := BeginRedisTx(ctx)
	for _, offset := range offsets {
		_, err := rtx.SetBit(rtx.ctx, b.key, int64(offset/redisBitSetMaxLength), 1).Result()
		if err != nil {
			rtx.Discard()
			return err
		}
	}
	_, err := rtx.Exec(rtx.ctx)
	return err
}

func (b *redisBitSet) Test(ctx context.Context, offsets []uint) (bool, error) {
	for _, offset := range offsets {
		res, err := b.rdb.GetBit(ctx, b.key, int64(offset/redisBitSetMaxLength)).Result()
		if err != nil || res == 0 {
			return false, err
		}
	}
	return true, nil
}
