package models

import (
	"context"
	"fmt"
	"time"
)

var keyTicketUsageCount = "ticket-count-%s"

func (tx rtx) IncrTicketUsageCount(ticketID string, incrCount int) (int64, error) {
	key := fmt.Sprintf(keyTicketUsageCount, ticketID)
	return tx.IncrBy(tx.ctx, key, int64(incrCount)).Result()
}

func InitTicketUsageCount(ctx context.Context, ticketID string, ttl time.Duration) error {
	key := fmt.Sprintf(keyTicketUsageCount, ticketID)
	return rdb.Set(ctx, key, 0, ttl).Err()
}
