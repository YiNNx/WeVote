package ticket

import (
	"github.com/YiNNx/WeVote/internal/common/errors"
	"github.com/YiNNx/WeVote/internal/pkg/ticket"
)

func ParseAndVerify(ticketStr string) (ticketID string, err error) {
	claims, err := ticket.ParseAndVerify(ticketStr)
	if err != nil {
		return "", errors.ErrInvalidTicket.WithErrDetail(err)
	}
	ticketID = claims.SubjectID
	return ticketID, nil
}
