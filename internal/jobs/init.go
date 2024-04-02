package jobs

import (
	"fmt"

	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/services"
	"github.com/YiNNx/WeVote/internal/utils/log"
	"github.com/YiNNx/WeVote/internal/utils/ticket"
	"github.com/robfig/cron/v3"
)

type Job struct {
	Spec string
	Func func()
}

var jobs = []Job{
	{
		Spec: fmt.Sprintf("@every %ds", config.C.Ticket.Spec),
		Func: func() {
			var err error
			services.Ticket, err = ticket.New()
			if err != nil {
				panic(err)
			}
		},
	},
}

func Init() {
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
