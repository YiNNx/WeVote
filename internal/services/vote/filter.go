package vote

import (
	"context"
	"time"

	"github.com/YiNNx/WeVote/internal/common/errors"
	"github.com/YiNNx/WeVote/internal/common/log"
	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/models"
	"github.com/YiNNx/WeVote/pkg/bloomfilter"
)

var filterUsername bloomfilter.BloomFilter

type UsernameSet map[string]struct{}

// DedupicateAndBloomFilter 合并重复参数，并使用布隆过滤器来检验是否含有无效字段
// 若含有无效字段则返回错误
func DedupicateAndBloomFilter(ctx context.Context, usernames []string) (UsernameSet, error) {
	if int64(len(usernames)) > config.C.Ticket.UpperLimit {
		return nil, errors.ErrTicketUsageExceed
	}

	usernameSet := make(UsernameSet, len(usernames))
	for _, username := range usernames {
		usernameSet[username] = struct{}{}
		err := ProcessBloomFilter(ctx, username)
		if err != nil {
			return nil, err
		}
	}
	return usernameSet, nil
}

// ProcessBloomFilter
func ProcessBloomFilter(ctx context.Context, username string) error {
	ok, err := filterUsername.Exists(ctx, []byte(username))
	if err != nil {
		return errors.ErrServerInternal.WithErrDetail(err)
	}
	if !ok {
		return errors.ErrInvalidUsernameExisted
	}
	return nil
}

func init() {
	filterUsername = bloomfilter.NewWithEstimates(
		1000000, 0.01,
		models.NewRedisBitSet("bitset-bloomfilter-username"),
	)

	usernames, err := models.GetAllUsernames()
	if err != nil {
		log.Logger.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	for _, username := range usernames {
		err := filterUsername.Add(ctx, []byte(username))
		if err != nil {
			log.Logger.Fatal(err)
		}
	}
}
