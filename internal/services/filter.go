package services

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

// DedupicateAndBloomFilter merges duplicate params and uses a bloom filter to check whether there are invalid usernames
// If it contains invalid fields, an error will be returned.
func DedupicateAndBloomFilter(ctx context.Context, usernames []string) (UsernameSet, error) {
	if len(usernames) > config.C.Ticket.UsageLimit {
		return nil, errors.ErrTicketLimitExceeded
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
		return errors.ErrInvalidUsername
	}
	return nil
}

func initBloomFilter() {
	filterUsername = bloomfilter.NewWithEstimates(
		100000, 0.0001,
		models.NewUserBloomfilterBitSet(),
	)

	usernames, err := models.FindAllExistedUser()
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
