package schema

import (
	"context"

	"github.com/golang-jwt/jwt"

	"github.com/YiNNx/WeVote/internal/common/errors"
	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/services"
	parser "github.com/YiNNx/WeVote/internal/utils/ticket"
	"github.com/YiNNx/WeVote/pkg/captcha"
)

const (
	statusFailed    = "failed"
	statusSucceeded = "succeeded"
)

// Vote is the resolver for the vote field.
func (r *mutationResolver) Vote(ctx context.Context, users []string, ticket string, recaptchaToken *string) (string, error) {
	if config.C.Captcha.Open {
		if recaptchaToken == nil {
			return statusFailed, errors.CaptchaTokenRequired
		}
		success, err := captcha.Client.Verify(*recaptchaToken)
		if err != nil || !success {
			return statusFailed, errors.CaptchaTokenInvalid
		}
	}

	if int64(len(users)) > config.C.Ticket.UpperLimit {
		return statusFailed, errors.TicketUsageLimitExceed
	}

	token, err := jwt.ParseWithClaims(ticket, &parser.TicketClaims{}, func(t *jwt.Token) (interface{}, error) { return config.C.Ticket.Secret, nil })
	if err != nil {
		return statusFailed, errors.TicketInvalid.WithErrDetail(err)
	}
	claims, ok := token.Claims.(*parser.TicketClaims)
	if !ok {
		return statusFailed, errors.TicketInvalid
	}
	ticketID := claims.SubjectId

	// remove duplicate data
	userSet := make(services.UserSet, len(users))
	for _, user := range users {
		userSet[user] = struct{}{}
	}

	err = services.Vote(ctx, ticketID, userSet)
	if err != nil {
		return statusFailed, err
	}
	return statusSucceeded, nil
}

// GetUserVotes is the resolver for the getUserVotes field.
func (r *queryResolver) GetUserVotes(ctx context.Context, username string) (*int, error) {
	count, err := services.GetVoteCount(ctx, username)
	return &count, err
}
