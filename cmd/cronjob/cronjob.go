package main

import (
	"os"

	"github.com/YiNNx/WeVote/internal/common/log"
	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/models"
	"github.com/YiNNx/WeVote/internal/services"
)

func main() {
	confPath := os.Getenv("CONF_PATH")
	config.Init(confPath)
	log.InitLogger()
	models.InitRedisClusterConns()

	services.StartCronJobTicketRefresh()

	for {
	}
}
