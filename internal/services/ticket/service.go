package ticket

import "github.com/YiNNx/WeVote/internal/config"

var (
	Parser       TicketParser
	GlobalTicket Ticket
)

func InitTicketService() {
	ticketProvider := newTicketProvider(
		config.C.Ticket.Secret,
		config.C.Ticket.Expiration.Duration,
	)
	ticketParser := newParser(config.C.Ticket.Secret)
	globalTicket := initGlobalTicket(ticketProvider)

	Parser = ticketParser
	GlobalTicket = globalTicket
}
