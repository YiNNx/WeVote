package models

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/YiNNx/WeVote/pkg/bloomfilter"
)


const (
	keyBloomfilterUsername = "bitset-bloomfilter-username"
)

func NewUserBloomfilterBitSet() bloomfilter.BitSetProvider {
	return newRedisBitSet(keyBloomfilterUsername)
}

type redisBitSet struct {
	rdb redis.UniversalClient
	key string
}

func newRedisBitSet(key string) bloomfilter.BitSetProvider {
	return &redisBitSet{rdb, key}
}

func (b *redisBitSet) Set(ctx context.Context, offsets []uint) error {
	rtx := b.rdb.Pipeline()
	for _, offset := range offsets {
		rtx.SetBit(ctx, b.key, int64(offset), 1)
	}
	_, err := rtx.Exec(ctx)
	return err
}

func (b *redisBitSet) Test(ctx context.Context, offsets []uint) (bool, error) {
	for _, offset := range offsets {
		res, err := b.rdb.GetBit(ctx, b.key, int64(offset)).Result()
		if err != nil || res == 0 {
			return false, err
		}
	}
	return true, nil
}
