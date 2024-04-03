package schema

import (
	"context"
)

// GetTicket is the resolver for the getTicket field.
func (r *queryResolver) GetTicket(ctx context.Context) (*string, error) {
	// return services.GetCurrentTicket()
	return nil, nil
}
