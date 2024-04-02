package jobs

import (
	"fmt"

	"github.com/robfig/cron/v3"

	"github.com/YiNNx/WeVote/pkg/log"
)

type Job struct {
	Spec string
	Func func()
}

func InitJobs(spec int) {
	jobs := []Job{
		{
			Spec: fmt.Sprintf("@every %ds", spec),
			Func: ticketGenerate,
		},
	}

	c := cron.New(cron.WithSeconds())

	for _, job := range jobs {
		_, err := c.AddFunc(job.Spec, job.Func)
		if err != nil {
			log.Logger.Error(err)
			return
		}
	}

	c.Start()
}
