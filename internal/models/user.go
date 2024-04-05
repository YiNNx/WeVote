package models

import (
	"time"

	"github.com/YiNNx/WeVote/pkg/bloomfilter"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"unique;not null"`
	VoteCount int    `gorm:"not null;default:0"`
}

func FindAllExistedUser() ([]string, error) {
	var usernames []string
	err := db.Model(&User{}).Select("username").Find(&usernames).Error
	return usernames, err
}

const (
	keyPrefixVote          = "vote-user:"
	keyPrefixTicketUsage   = "usage-ticket:"
	keyBloomfilterUsername = "bitset-bloomfilter-username"
)

func NewUserBloomfilterBitSet() bloomfilter.BitSetProvider {
	return newRedisBitSet(keyBloomfilterUsername)
}

func NewTicketUsageCounter(ttl time.Duration, limit int) *Counter {
	return &Counter{
		rdb:       rdb,
		keyPrefix: keyPrefixTicketUsage,
		ttl:       ttl,
		limit:     limit,
	}
}
