package schema

import (
	"context"
	"fmt"
)

// Vote is the resolver for the vote field.
func (r *mutationResolver) Vote(ctx context.Context, users []string, ticket string) (*string, error) {
	panic(fmt.Errorf("not implemented: Vote - vote"))
}

// GetUserVotes is the resolver for the getUserVotes field.
func (r *queryResolver) GetUserVotes(ctx context.Context, username string) (*int, error) {
	panic(fmt.Errorf("not implemented: GetUserVotes - getUserVotes"))
}
