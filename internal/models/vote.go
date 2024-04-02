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
var keyVoteCountModifiedSet = "vote-count-modified"

func (tx rtx) IncrVoteCount(user string) (count int64, err error) {
	key := fmt.Sprintf(keyTicketUsageCount, user)
	return tx.Incr(tx.ctx, key).Result()
}

func (tx rtx) RecordVoteCountModified(user []string) error {
	return tx.SAdd(tx.ctx, keyVoteCountModifiedSet, user).Err()
}

func GetVoteCount(ctx context.Context, user string) (count string, err error) {
	key := fmt.Sprintf(keyTicketUsageCount, user)
	res := rdb.Get(ctx, key)
	return res.Val(), res.Err()
}
