package services

import "github.com/YiNNx/WeVote/internal/services/cron"

func InitServices() {
	cron.InitJobs()
}
