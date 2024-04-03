package models

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"unique;not null"`
	VoteCount int    `gorm:"not null;default:0"`
}

func GetAllUsernames() ([]string, error) {
	var usernames []string
	err := db.Model(&User{}).Select("username").Find(&usernames).Error
	return usernames, err
}

var keyVoteCount = "vote-count-%s"
var keyVoteCountModifiedSet = "vote-count-modified"

func (tx rtx) IncrVoteCount(user string) (count int64, err error) {
	key := fmt.Sprintf(keyVoteCount, user)
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
