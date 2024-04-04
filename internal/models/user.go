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

func GetVoteCountByUsername(username string) (int, error) {
	var user User
	err := db.Where("username = ?", username).Take(&user).Error
	return user.VoteCount, err
}

type VoteCount2Usernames map[int][]string

func (tx *tx) UpdateVoteCountBatch(votes VoteCount2Usernames) error {
	for count, usernames := range votes {
		err := tx.Model(&User{}).
			Where("username IN ?", usernames).
			Update("vote_count", count).Error
		if err != nil {
			return err
		}
	}
	return nil
}

var keyVoteCount = "vote-count-%s"
var keySetUserUpdated = "vote-count-modified"

func (tx rtx) IncrVoteCount(user string) (count int64, err error) {
	key := fmt.Sprintf(keyVoteCount, user)
	return tx.Incr(tx.ctx, key).Result()
}

func (tx rtx) SetVoteCount(user string, count int) (err error) {
	key := fmt.Sprintf(keyVoteCount, user)
	return tx.Set(tx.ctx, key, count, 0).Err()
}

func (tx rtx) RecordUserModified(user []string) error {
	return tx.SAdd(tx.ctx, keySetUserUpdated, user).Err()
}

func GetUsersModified(ctx context.Context) ([]string, error) {
	res := rdb.SMembers(ctx, keySetUserUpdated)
	return res.Val(), res.Err()
}

func GetVoteCountByCache(ctx context.Context, user string) (count string, err error) {
	key := fmt.Sprintf(keyTicketUsageCount, user)
	res := rdb.Get(ctx, key)
	return res.Val(), res.Err()
}
