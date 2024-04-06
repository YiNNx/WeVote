package services

import (
	"context"

	"github.com/YiNNx/WeVote/internal/common/errors"
	"github.com/YiNNx/WeVote/internal/common/log"
	"github.com/YiNNx/WeVote/internal/models"
)

// For illegal requests: If intercepted by the Bloom filter (success rate >= 99.99%), the entire batch will not be executed.
// If it is not intercepted, it will be executed in sequence. If the execution of illegal fields fails, it will be skipped.
func ProcessVote(ctx context.Context, ticketID string, users UsernameSet) error {
	for user := range users {
		err := models.VoteDataWrapper.Incr(ctx, user)
		if err != nil && err != errors.ErrInvalidKey {
			log.Logger.Error(err)
		}
	}
	return nil
}

func GetVoteCount(ctx context.Context, user string) (*int, error) {
	return models.VoteDataWrapper.Get(ctx, user)
}
