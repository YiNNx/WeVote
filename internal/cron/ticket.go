package cron

import (
	"github.com/YiNNx/WeVote/internal/services"
	"github.com/YiNNx/WeVote/internal/utils/ticket"
)

func ticketGenerate() {
	var err error
	services.Ticket, err = ticket.GenerateTicket()
	if err != nil {
		panic(err)
	}
}
