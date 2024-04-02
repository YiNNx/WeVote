package cron

import (
	"context"

	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/models"
	"github.com/YiNNx/WeVote/internal/services"
	"github.com/YiNNx/WeVote/internal/utils/ticket"
	"github.com/YiNNx/WeVote/pkg/log"
)

func jobTicketGenerate() {
	ticketID, ticketStr, err := ticket.GenerateTicket()
	if err != nil {
		log.Logger.Error(err)
	}
	services.Ticket = &ticketStr
	err = models.InitTicketUsageCount(context.Background(), ticketID, config.C.Ticket.Expiration.Duration)
	if err != nil {
		log.Logger.Error(err)
	}
}
