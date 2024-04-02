package schema

import (
	"context"
	"fmt"

	"github.com/YiNNx/WeVote/internal/common/errors"
	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/services"
)

// Vote is the resolver for the vote field.
func (r *mutationResolver) Vote(ctx context.Context, users []string, ticket string) (*string, error) {
	if int64(len(users)) > config.C.Ticket.UpperLimit {
		return nil, errors.TicketUsageLimitExceed
	}

	// remove duplicate data
	userSet := make(services.UserSet, len(users))
	for _, user := range users {
		userSet[user] = struct{}{}
	}
	panic(fmt.Errorf("not implemented: Vote - vote"))
}

// GetUserVotes is the resolver for the getUserVotes field.
func (r *queryResolver) GetUserVotes(ctx context.Context, username string) (*int, error) {
	panic(fmt.Errorf("not implemented: GetUserVotes - getUserVotes"))
}
