package schema

import (
	"context"

	"github.com/YiNNx/WeVote/internal/services"
)

const (
	statusFailed    = "failed"
	statusSucceeded = "succeeded"
)

// Vote is the resolver for the vote field.
func (r *mutationResolver) Vote(ctx context.Context, users []string, ticket string, recaptchaToken *string) (bool, error) {
	err := services.VerifyCaptchaIfCaptchaOpened(recaptchaToken)
	if err != nil {
		return false, err
	}

	ticketID, err := services.ParseAndVerifyTicket(ticket)
	if err != nil {
		return false, err
	}

	userSet, err := services.DedupicateAndBloomFilter(ctx, users)
	if err != nil {
		return false, err
	}

	err = services.ProcessVote(ctx, ticketID, userSet)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetUserVotes is the resolver for the getUserVotes field.
func (r *queryResolver) GetUserVotes(ctx context.Context, username string) (*int, error) {
	err := services.ProcessBloomFilter(ctx, username)
	if err != nil {
		return nil, err
	}
	return services.GetVoteCount(ctx, username)
}
