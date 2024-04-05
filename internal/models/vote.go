package models

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/pkg/cache"
)

var VoteDataWrapper cache.WriteBehindWrapper[string, int]

func initVoteDataWrapper() {
	VoteDataWrapper = cache.NewWriteBehindDataWrapper(
		cache.WriteBehindConfig{WriteBackPeriod: config.C.Vote.WriteBackPeriod.Duration},
		newVoteCacheWrapper(rdb, keyPrefixVote, config.C.Vote.CacheTTL),
		newVoteDBWrapper(db),
		newVoteDirtyKeyTracker(),
	)
}

type voteCacheWrapper struct {
	redis.UniversalClient
	keyPrefix string
	ttl       time.Duration
}

func (cache *voteCacheWrapper) Key(username string) string {
	return cache.keyPrefix + username
}

func (cache *voteCacheWrapper) GetTTL() time.Duration {
	return cache.ttl
}

func (cache *voteCacheWrapper) CacheGet(ctx context.Context, username string) (*int, error) {
	res, err := cache.Get(ctx, cache.Key(username)).Int()
	if err == redis.Nil {
		return nil, nil
	}
	return &res, err
}

func (cache *voteCacheWrapper) CacheGetBatch(ctx context.Context, usernames []string) (userVotes map[string]int, err error) {
	userVotes = make(map[string]int, len(usernames))

	pipe := cache.Pipeline()
	results := make([]*redis.StringCmd, len(usernames))

	for i := range usernames {
		results[i] = pipe.Get(ctx, cache.Key(usernames[i]))
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	for i, res := range results {
		val, err := res.Int()
		if err == nil {
			userVotes[usernames[i]] = val
		}
	}
	return userVotes, nil
}

func (cache *voteCacheWrapper) CacheIncr(ctx context.Context, username string) error {
	cacheKey := cache.Key(username)
	err := cache.Incr(ctx, cacheKey).Err()
	if err != nil {
		return err
	}
	return cache.Expire(ctx, cacheKey, cache.GetTTL()).Err()
}

func newVoteCacheWrapper(rdb redis.UniversalClient, keyPrefix string, ttl time.Duration) *voteCacheWrapper {
	return &voteCacheWrapper{
		UniversalClient: rdb,
		keyPrefix:       keyPrefix,
		ttl:             ttl,
	}
}

type voteDBWrapper struct {
	*gorm.DB
}

func (db *voteDBWrapper) Select(username string) (*int, error) {
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user.VoteCount, err
}

func (db *voteDBWrapper) UpdateBatch(userVotes map[string]int) error {
	vote2users := make(map[int][]string)

	for user, vote := range userVotes {
		if users, ok := vote2users[vote]; ok {
			vote2users[vote] = append(users, user)
		} else {
			vote2users[vote] = []string{user}
		}
	}

	for vote, users := range vote2users {
		err := db.Model(&User{}).
			Where("username IN ?", users).
			Update("vote_count", vote).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func newVoteDBWrapper(db *gorm.DB) *voteDBWrapper {
	return &voteDBWrapper{db}
}

type voteDirtyKeyTracker struct {
	sync.RWMutex
	usernames []string
}

func (tracker *voteDirtyKeyTracker) GetAll() []string {
	return tracker.usernames
}

func (tracker *voteDirtyKeyTracker) Put(username string) {
	tracker.usernames = append(tracker.usernames, username)
}

func (tracker *voteDirtyKeyTracker) Clear() {
	tracker.usernames = []string{}
}

func newVoteDirtyKeyTracker() *voteDirtyKeyTracker {
	return &voteDirtyKeyTracker{
		RWMutex:   sync.RWMutex{},
		usernames: []string{},
	}
}
