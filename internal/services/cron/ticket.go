package cron

import (
	"context"
	"fmt"

	"github.com/YiNNx/WeVote/internal/common/log"
	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/models"
	"github.com/YiNNx/WeVote/internal/services/ticket"
)

var cronJobTicketGrant = cronJob{
	Spec: fmt.Sprintf("@every %s", config.C.Ticket.Expiration.Duration),
	Func: ticketGrant,
}

func ticketGrant() {
	ticketID, err := ticket.GlobalTicket.Grant()
	if err != nil {
		log.Logger.Error(err)
	}
	// services.SharedTicket = &ticketStr
	err = models.InitTicketUsageCount(context.Background(), ticketID, config.C.Ticket.Expiration.Duration)
	if err != nil {
		log.Logger.Error(err)
	}
}
