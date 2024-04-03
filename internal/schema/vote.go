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
func (r *mutationResolver) Vote(ctx context.Context, users []string, ticket string, recaptchaToken *string) (string, error) {
	err := services.VerifyCaptchaIfCaptchaOpen(recaptchaToken)
	if err != nil {
		return statusFailed, err
	}

	ticketID, err := services.ParseAndVerifyTicket(ticket)
	if err != nil {
		return statusFailed, err
	}

	userSet, err := services.DedupicateAndBloomFilter(ctx, users)
	if err != nil {
		return statusFailed, err
	}

	err = services.Vote(ctx, ticketID, userSet)
	if err != nil {
		return statusFailed, err
	}

	return statusSucceeded, nil
}

// GetUserVotes is the resolver for the getUserVotes field.
func (r *queryResolver) GetUserVotes(ctx context.Context, username string) (*int, error) {
	err := services.BloomFilter(ctx, username)
	if err != nil {
		return nil, err
	}
	count, err := services.GetVoteCount(ctx, username)
	return &count, err
}
