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

const (
	keyPrefixVote = "vote-user:"
)

type Vote struct {
	gorm.Model
	Username  string `gorm:"unique;not null"`
	VoteCount int    `gorm:"not null;default:0"`
}

func FindAllExistedUser() ([]string, error) {
	var usernames []string
	err := db.Model(&Vote{}).Select("username").Find(&usernames).Error
	return usernames, err
}

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
	rdb       redis.UniversalClient
	keyPrefix string
	ttl       time.Duration
}

func (c *voteCacheWrapper) key(username string) string {
	return c.keyPrefix + username
}

func (c *voteCacheWrapper) Get(ctx context.Context, username string) (*int, error) {
	res, err := c.rdb.Get(ctx, c.key(username)).Int()
	if err == redis.Nil {
		return nil, nil
	}
	return &res, err
}

func (c *voteCacheWrapper) SetNX(ctx context.Context, username string, vote int) (bool, error) {
	return c.rdb.SetNX(ctx, c.key(username), vote, c.ttl).Result()
}

func (c *voteCacheWrapper) GetBatch(ctx context.Context, usernames []string) (userVotes map[string]int, err error) {
	userVotes = make(map[string]int, len(usernames))

	pipe := c.rdb.Pipeline()
	results := make([]*redis.StringCmd, len(usernames))

	for i := range usernames {
		results[i] = pipe.Get(ctx, c.key(usernames[i]))
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

func (c *voteCacheWrapper) Incr(ctx context.Context, username string) error {
	cacheKey := c.key(username)
	err := c.rdb.Incr(ctx, cacheKey).Err()
	if err != nil {
		return err
	}
	return c.rdb.Expire(ctx, cacheKey, c.ttl).Err()
}

func newVoteCacheWrapper(rdb redis.UniversalClient, keyPrefix string, ttl time.Duration) *voteCacheWrapper {
	return &voteCacheWrapper{
		rdb:       rdb,
		keyPrefix: keyPrefix,
		ttl:       ttl,
	}
}

type voteDBWrapper struct {
	*gorm.DB
}

func (db *voteDBWrapper) Select(username string) (*int, error) {
	var user Vote
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
		err := db.Model(&Vote{}).
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
	sync.Mutex
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
		Mutex:     sync.Mutex{},
		usernames: []string{},
	}
}
