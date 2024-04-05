package ticket

import (
	"context"
	"fmt"

	"github.com/YiNNx/WeVote/internal/pkg/cron"
	"github.com/YiNNx/WeVote/internal/common/log"
	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/models"
)

// cronJobTicketGrant periodically issues and updates the global ticket
var cronJobTicketGrant = cron.CronJob{
	Spec: fmt.Sprintf("@every %s", config.C.Ticket.Expiration.Duration),
	Func: ticketGrant,
}

func ticketGrant() {
	ticketID, err := globalGrant()
	if err != nil {
		log.Logger.Error(err)
	}
	err = models.StartTicketUsageCount(context.Background(), ticketID, config.C.Ticket.Expiration.Duration)
	if err != nil {
		log.Logger.Error(err)
	}
}

func init() {
	cronJobTicketGrant.Start()
}
