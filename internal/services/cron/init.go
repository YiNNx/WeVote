package cron

import (
	"github.com/robfig/cron/v3"

	"github.com/YiNNx/WeVote/internal/common/log"
)

type cronJob struct {
	Spec string
	Func func()
}

var cronJobs = []cronJob{
	cronJobTicketGrant,
}

func InitJobs() {
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
