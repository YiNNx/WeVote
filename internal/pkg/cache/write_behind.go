package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/YiNNx/WeVote/internal/common/errors"
	"github.com/YiNNx/WeVote/internal/common/log"
	"github.com/YiNNx/WeVote/internal/pkg/cron"
)

type WriteBehindWrapper[K, V comparable] interface {
	Get(ctx context.Context, key K) (val *V, err error)
	Incr(ctx context.Context, key K) error
	WriteBack(ctx context.Context) error
}

type writeBehindWrapper[K, V comparable] struct {
	cache           cacheRedis[K, V]
	database        database[K, V]
	dirtyKeyTracker dirtyKeyTracker[K]
}

type cacheRedis[K, V comparable] interface {
	redis.UniversalClient

	Key(K) string
	GetTTL() time.Duration
	CacheGet(context.Context, K) (*V, error)
	CacheGetBatch(context.Context, []K) (map[K]V, error)
	CacheIncr(context.Context, K) error
}

type database[K, V comparable] interface {
	Select(K) (*V, error)
	UpdateBatch(map[K]V) error
}

type dirtyKeyTracker[K comparable] interface {
	RWMutex

	GetAll() []K
	Put(K)
	Clear()
}

type RWMutex interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()
}

func (p *writeBehindWrapper[K, V]) Get(ctx context.Context, key K) (val *V, err error) {
	// If hit the cache, return the val directly
	val, err = p.cache.CacheGet(ctx, key)
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

	// If the val exists in the db, use the optimistic lock (Redis WATCH) to update cache
	keyWatch := p.cache.Key(key)
	txWatch := func(tx *redis.Tx) error {
		// If the key has been or is being updated by other threads, the below operations will fail
		cacheExists, err := tx.Exists(ctx, keyWatch).Result()
		if err != nil {
			return err
		}
		if cacheExists != 0 {
			return redis.TxFailedErr
		}
		pipe := tx.TxPipeline()
		err = pipe.Set(ctx, keyWatch, val, p.cache.GetTTL()).Err()
		if err != nil {
			return err
		}
		_, err = pipe.Exec(ctx)
		return err
	}
	err = p.cache.Watch(ctx, txWatch, keyWatch)

	// if the err is caused by the optimistic lock, just re-read the cache and return
	if err == redis.TxFailedErr {
		return p.cache.CacheGet(ctx, key)
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

	// cause the incr operation is atomic, there's no need to use lock
	// flag the key dirty and incr
	p.dirtyKeyTracker.RLock()
	defer p.dirtyKeyTracker.RUnlock()

	p.dirtyKeyTracker.Put(key)

	return p.cache.CacheIncr(ctx, key)
}

func (p *writeBehindWrapper[K, V]) WriteBack(ctx context.Context) error {
	// Lock dirtyTracker to avoid data race during writeback
	p.dirtyKeyTracker.Lock()
	defer p.dirtyKeyTracker.Unlock()

	dirtyKeys := p.dirtyKeyTracker.GetAll()

	dirtyData, err := p.cache.CacheGetBatch(ctx, dirtyKeys)
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

func NewWriteBehindDataWrapper[K, V comparable](config WriteBehindConfig, cache cacheRedis[K, V], database database[K, V], dirtyKeyTracker dirtyKeyTracker[K]) WriteBehindWrapper[K, V] {
	wrapper := &writeBehindWrapper[K, V]{
		cache:           cache,
		database:        database,
		dirtyKeyTracker: dirtyKeyTracker,
	}
	writeBackJob := func() {
		wrapper.WriteBack(context.Background())
		log.Logger.Info("proccess write-back")
	}
	cron.NewCronJob(config.WriteBackPeriod, writeBackJob).Start()
	return wrapper
}
