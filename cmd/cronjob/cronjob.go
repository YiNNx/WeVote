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
	err := config.Init(confPath)
	if err != nil {
		log.Logger.Fatal(err)
	}
	log.InitLogger()
	models.InitIOWrapper()

	services.StartCronJobTicketRefresh()

	for {
	}
}
