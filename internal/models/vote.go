package models

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/YiNNx/WeVote/internal/pkg/cache"
)

var voteDataWrapper cache.WriteBehindWrapper[string, int]

func initVoteDataWrapper(rdb redis.UniversalClient, db *gorm.DB) {
	voteDataWrapper = cache.NewWriteBehindDataWrapper(
		newVoteCacheWrapper(rdb),
		newVoteDBWrapper(db),
		newVoteDirtyKeyTracker(),
	)
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

type voteCacheWrapper struct {
	redis.UniversalClient
}

func (cache *voteCacheWrapper) Key(username string) string {
	return fmt.Sprintf("vote-count-%s", username)
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
	return cache.Incr(ctx, cache.Key(username)).Err()
}

func newVoteCacheWrapper(rdb redis.UniversalClient) *voteCacheWrapper {
	return &voteCacheWrapper{rdb}
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

// var dataWrapper *cache.WriteBehindDataWrapper

// var cronJobWriteBehind = cron.CronJob{
// 	Spec: fmt.Sprintf("@every %s", "1s"),
// 	Func: writeBehindJob,
// }

// func writeBehindJob() {
// 	// ctx := context.Background()
// 	// err := dataWrapper.WriteBack(ctx)
// 	// if err != nil {
// 	// 	log.Logger.Error(err)
// 	// }
// }

// func init() {
// 	cronJobWriteBehind.Start()
// }

// func writeBehindJob() {
// 	ctx := context.Background()
// 	usernames, err := GetUsersModified(ctx)
// 	if err != nil {
// 		log.Logger.Error(err)
// 	}
// 	votes := make(VoteCount2Usernames)
// 	for _, username := range usernames {
// 		res, err := GetVoteCountByCache(ctx, username)
// 		if err != nil {
// 			log.Logger.Error(err)
// 			continue
// 		}
// 		count, err := strconv.Atoi(res)
// 		if err != nil {
// 			log.Logger.Error(err)
// 			continue
// 		}
// 		if users, ok := votes[count]; ok {
// 			votes[count] = append(users, username)
// 		} else {
// 			votes[count] = []string{username}
// 		}
// 	}
// 	tx := BeginPostgresTx()
// 	err = tx.UpdateVoteCountBatch(votes)
// 	if err != nil {
// 		tx.Rollback()
// 		log.Logger.Error(err)
// 	}
// 	tx.Commit()
// }
