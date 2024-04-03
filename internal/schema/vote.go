package schema

import (
	"context"

	"github.com/YiNNx/WeVote/internal/services/captcha"
	ticketsrv "github.com/YiNNx/WeVote/internal/services/ticket"
	"github.com/YiNNx/WeVote/internal/services/vote"
)

const (
	statusFailed    = "failed"
	statusSucceeded = "succeeded"
)

// Vote is the resolver for the vote field.
func (r *mutationResolver) Vote(ctx context.Context, users []string, ticket string, recaptchaToken *string) (string, error) {
	err := captcha.VerifyCaptchaIfCaptchaOpen(recaptchaToken)
	if err != nil {
		return statusFailed, err
	}

	ticketClaims, err := ticketsrv.Parser.ParseAndVerify(ticket)
	if err != nil {
		return statusFailed, err
	}

	userSet, err := vote.DedupicateAndBloomFilter(ctx, users)
	if err != nil {
		return statusFailed, err
	}

	err = vote.Vote(ctx, ticketClaims.SubjectId, userSet)
	if err != nil {
		return statusFailed, err
	}

	return statusSucceeded, nil
}

// GetUserVotes is the resolver for the getUserVotes field.
func (r *queryResolver) GetUserVotes(ctx context.Context, username string) (*int, error) {
	err := vote.BloomFilter(ctx, username)
	if err != nil {
		return nil, err
	}
	count, err := vote.GetVoteCount(ctx, username)
	return &count, err
}
