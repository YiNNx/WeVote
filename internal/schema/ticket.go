package schema

import (
	"context"

	"github.com/YiNNx/WeVote/internal/services"
)

// GetTicket is the resolver for the getTicket field.
func (r *queryResolver) GetTicket(ctx context.Context) (string, error) {
	return services.GetTicket(), nil
}
