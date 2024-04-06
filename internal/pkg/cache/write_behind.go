package cache

import (
	"context"
	"time"

	"github.com/YiNNx/WeVote/internal/common/errors"
	"github.com/YiNNx/WeVote/internal/common/log"
	"github.com/YiNNx/WeVote/pkg/cronjob"
)

type WriteBehindWrapper[K, V comparable] interface {
	Get(ctx context.Context, key K) (val *V, err error)
	Incr(ctx context.Context, key K) error
	WriteBack(ctx context.Context) error
}

type writeBehindWrapper[K, V comparable] struct {
	cache           cache[K, V]
	database        database[K, V]
	dirtyKeyTracker dirtyKeyTracker[K]
}

type cache[K, V comparable] interface {
	Get(context.Context, K) (*V, error)
	GetBatch(context.Context, []K) (map[K]V, error)
	Incr(context.Context, K) error
	SetNX(context.Context, K, V) (bool, error)
}

type database[K, V comparable] interface {
	Select(K) (*V, error)
	UpdateBatch(map[K]V) error
}

type dirtyKeyTracker[K comparable] interface {
	GetAll() []K
	Put(K)
	Clear()

	Mutex
}

type Mutex interface {
	Lock()
	Unlock()
}

func (p *writeBehindWrapper[K, V]) Get(ctx context.Context, key K) (val *V, err error) {
	// If hit the cache, return the val directly
	val, err = p.cache.Get(ctx, key)
	if err != nil || val != nil {
		return val, err
	}

	// If not hit the cache, check if the val exists in the database
	val, err = p.database.Select(key)
	if err != nil {
		return nil, err
	}
	// if the val not exists in the db, return
	if val == nil {
		return nil, nil
	}

	// Use the optimistic lock to update cache
	ok, err := p.cache.SetNX(ctx, key, *val)
	if err != nil {
		return nil, err
	}
	if !ok {
		return p.cache.Get(ctx, key)
	}

	return val, err
}

func (p *writeBehindWrapper[K, V]) Incr(ctx context.Context, key K) error {
	// hit cache or sync cache from db
	val, err := p.Get(ctx, key)
	if err != nil {
		return err
	}
	if val == nil {
		return errors.ErrInvalidKey
	}

	p.dirtyKeyTracker.Lock()
	defer p.dirtyKeyTracker.Unlock()

	p.dirtyKeyTracker.Put(key)

	return p.cache.Incr(ctx, key)
}

func (p *writeBehindWrapper[K, V]) WriteBack(ctx context.Context) error {
	// Lock dirtyTracker to avoid data race during writeback
	p.dirtyKeyTracker.Lock()
	defer p.dirtyKeyTracker.Unlock()

	dirtyKeys := p.dirtyKeyTracker.GetAll()

	dirtyData, err := p.cache.GetBatch(ctx, dirtyKeys)
	if err != nil {
		return err
	}

	err = p.database.UpdateBatch(dirtyData)
	if err != nil {
		return err
	}

	p.dirtyKeyTracker.Clear()

	return nil
}

type WriteBehindConfig struct {
	WriteBackPeriod time.Duration
}

func NewWriteBehindDataWrapper[K, V comparable](config WriteBehindConfig, cache cache[K, V], database database[K, V], dirtyKeyTracker dirtyKeyTracker[K]) WriteBehindWrapper[K, V] {
	wrapper := &writeBehindWrapper[K, V]{
		cache:           cache,
		database:        database,
		dirtyKeyTracker: dirtyKeyTracker,
	}
	writeBackJob := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		err := wrapper.WriteBack(ctx)
		if err != nil {
			log.Logger.Error(err)
		}
		log.Logger.Info("proccess write-back")
	}
	cronjob.NewCronJob(config.WriteBackPeriod, writeBackJob).Start()
	return wrapper
}
