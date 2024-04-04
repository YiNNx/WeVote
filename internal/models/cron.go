package models

import (
	"context"
	"fmt"
	"strconv"

	"github.com/robfig/cron/v3"

	"github.com/YiNNx/WeVote/internal/common/log"
)

type cronJob struct {
	Spec string
	Func func()
}

var cronJobs = []cronJob{
	cronJobWriteBehind,
}

var cronJobWriteBehind = cronJob{
	Spec: fmt.Sprintf("@every %s", "1s"),
	Func: writeBehind,
}

func writeBehind() {
	ctx := context.Background()
	usernames, err := GetUsersModified(ctx)
	if err != nil {
		log.Logger.Error(err)
	}
	votes := make(VoteCount2Usernames)
	for _, username := range usernames {
		res, err := GetVoteCountByCache(ctx, username)
		if err != nil {
			log.Logger.Error(err)
			continue
		}
		count, err := strconv.Atoi(res)
		if err != nil {
			log.Logger.Error(err)
			continue
		}
		if users, ok := votes[count]; ok {
			votes[count] = append(users, username)
		} else {
			votes[count] = []string{username}
		}
	}
	tx := BeginPostgresTx()
	err = tx.UpdateVoteCountBatch(votes)
	if err != nil {
		tx.Rollback()
		log.Logger.Error(err)
	}
	tx.Commit()
}

func InitCronJobs() {
	c := cron.New(cron.WithSeconds())

	for _, job := range cronJobs {
		_, err := c.AddFunc(job.Spec, job.Func)
		if err != nil {
			log.Logger.Error(err)
			return
		}
	}

	c.Start()
}
