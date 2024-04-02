package cron

import (
	"fmt"

	"github.com/robfig/cron/v3"

	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/pkg/log"
)

type CronJob struct {
	Spec string
	Func func()
}

var cronJobs = []CronJob{
	{
		Spec: fmt.Sprintf("@every %s", config.C.Ticket.Expiration.Duration),
		Func: ticketGenerate,
	},
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
