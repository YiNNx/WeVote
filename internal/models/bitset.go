package models

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const (
	keyVoteBloomFilter   = "vote-bloom-filter"
	redisBitSetMaxLength = 4 * 1024 * 1024 * 1024
)

type RedisBitSet struct {
	rdb redis.UniversalClient
	key string
}

func NewRedisBitSet() *RedisBitSet {
	return &RedisBitSet{rdb, keyVoteBloomFilter}
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
