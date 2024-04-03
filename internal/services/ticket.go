package services

import (
	"github.com/YiNNx/WeVote/internal/common/errors"
	"github.com/YiNNx/WeVote/internal/utils/ticket"
)

var sharedTicket *string

func GetCurrentTicket() (*string, error) {
	if sharedTicket == nil {
		return nil, errors.ErrGetTicket
	}
	return sharedTicket, nil
}

func ParseAndVerifyTicket(ticketStr string) (ticketID string, err error) {
	claims, err := ticket.ParseAndVerifyTicket(ticketStr)
	if err != nil {
		return "", errors.ErrInvalidTicket.WithErrDetail(err)
	}
	ticketID = claims.SubjectId
	return ticketID, nil
}
