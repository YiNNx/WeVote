package services

import "github.com/YiNNx/WeVote/internal/common/errors"

var Ticket *string

func GetTicket() (*string, error) {
	if Ticket == nil {
		return nil, errors.GetTicketFailed
	}
	return Ticket, nil
}
