package cron

import (
	"fmt"

	"github.com/YiNNx/WeVote/internal/config"
)

var cronJobTicketGrant = cronJob{
	Spec: fmt.Sprintf("@every %s", config.C.Ticket.Expiration.Duration),
	Func: ticketGrant,
}

func ticketGrant() {
	// ticketID, ticketStr, err := ticket.GenerateTicket()
	// if err != nil {
	// 	log.Logger.Error(err)
	// }
	// // services.SharedTicket = &ticketStr
	// err = models.InitTicketUsageCount(context.Background(), ticketID, config.C.Ticket.Expiration.Duration)
	// if err != nil {
	// 	log.Logger.Error(err)
	// }
}
