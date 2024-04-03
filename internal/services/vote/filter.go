package vote

import (
	"context"
	"time"

	"github.com/YiNNx/WeVote/internal/common/errors"
	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/models"
	"github.com/YiNNx/WeVote/pkg/bloomfilter"
)

var bloomFilterUsername bloomfilter.BloomFilter

type UsernameSet map[string]struct{}

// DedupicateAndBloomFilter 合并重复参数，并使用布隆过滤器来检验是否含有无效字段，若含有无效字段则返回错误
func DedupicateAndBloomFilter(ctx context.Context, usernames []string) (UsernameSet, error) {
	if int64(len(usernames)) > config.C.Ticket.UpperLimit {
		return nil, errors.ErrTicketUsageExceed
	}

	usernameSet := make(UsernameSet, len(usernames))
	for _, username := range usernames {
		usernameSet[username] = struct{}{}
		err := BloomFilter(ctx, username)
		if err != nil {
			return nil, err
		}
	}
	return usernameSet, nil
}

// BloomFilterBatch 合并重复参数，并使用布隆过滤器来检验是否含有无效字段，若含有无效字段则返回错误
func BloomFilter(ctx context.Context, username string) error {
	ok, err := bloomFilterUsername.Exists(ctx, []byte(username))
	if err != nil {
		return errors.ErrServerInternal.WithErrDetail(err)
	}
	if !ok {
		return errors.ErrInvalidUsernameExisted
	}
	return nil
}

func InitUsernameBloomFilter() error {
	bloomFilterUsername = bloomfilter.NewWithEstimates(
		1000000, 0.01,
		models.NewRedisBitSet(),
	)

	usernames, err := models.GetAllUsernames()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	for _, username := range usernames {
		err := bloomFilterUsername.Add(ctx, []byte(username))
		if err != nil {
			return err
		}
	}
	return nil
}
