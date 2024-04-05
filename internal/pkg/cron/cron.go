package cron

import (
	"github.com/robfig/cron/v3"

	"github.com/YiNNx/WeVote/internal/common/log"
)

type CronJob struct {
	Spec string
	Func func()
}

func (job *CronJob) Start() {
	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc(job.Spec, job.Func)
	if err != nil {
		log.Logger.Error(err)
		return
	}

	c.Start()
}
