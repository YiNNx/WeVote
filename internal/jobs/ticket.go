package jobs

import (
	"github.com/YiNNx/WeVote/internal/services"
	"github.com/YiNNx/WeVote/internal/utils/ticket"
)

func ticketGenerate() {
	var err error
	services.Ticket, err = ticket.Generator.Generate()
	if err != nil {
		panic(err)
	}
}
