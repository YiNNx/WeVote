package services

import (
	"sync"

	"github.com/YiNNx/WeVote/internal/common/errors"
	"github.com/YiNNx/WeVote/internal/common/ticket"
)

type SharedTicket struct {
	ticket   string
	provider ticket.Provider
	mutex    sync.RWMutex
}

func (t *SharedTicket) Grant() (tid string, err error) {
	tid, ticket, err := t.provider.Generate()
	if err != nil {
		return "", err
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.ticket = ticket

	return tid, nil
}

func (t *SharedTicket) Access() string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.ticket
}

func InitSharedTicket(provider ticket.Provider) (*SharedTicket, error) {
	_, initTicket, err := provider.Generate()
	if err != nil {
		return nil, err
	}
	return &SharedTicket{
		ticket:   initTicket,
		provider: provider,
		mutex:    sync.RWMutex{},
	}, nil
}

func ParseAndVerifyTicket(ticketStr string) (ticketID string, err error) {
	claims, err := ticket.ParseAndVerifyTicket(ticketStr)
	if err != nil {
		return "", errors.ErrInvalidTicket.WithErrDetail(err)
	}
	ticketID = claims.SubjectId
	return ticketID, nil
}
