package models

import (
	"context"
	"fmt"
	"time"
)

var keyTicketUsageCount = "ticket-count-%s"

func (tx rtx) IncrTicketUsageCount(ticketID string, incrCount int) (int64, error) {
	key := fmt.Sprintf(keyTicketUsageCount, ticketID)
	ctx := context.Background()
	return tx.IncrBy(ctx, key, int64(incrCount)).Result()
}


func InitializeKeyUsageCount(ticketID string, ttl time.Duration) error {
	key := fmt.Sprintf(keyTicketUsageCount, ticketID)
	ctx := context.Background()
	return rdb.Set(ctx, key, 0, ttl).Err()
}
