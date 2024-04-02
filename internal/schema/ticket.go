package schema

import (
	"context"
	"errors"
)

// GetTicket is the resolver for the getTicket field.
func (r *queryResolver) GetTicket(ctx context.Context) (*string, error) {
	return nil, errors.New("err msg")
}
