package models

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Vote struct {
	gorm.Model
	Username string
	Count    int
}

var keyVoteCount = "vote-count-%s"

func (tx rtx) IncrVoteCount(username string) (count int64, err error) {
	ctx := context.Background()
	key := fmt.Sprintf(keyTicketUsageCount, username)
	return tx.Incr(ctx, key).Result()
}
