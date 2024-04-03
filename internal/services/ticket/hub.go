package ticket

import (
	"sync"

	"github.com/YiNNx/WeVote/internal/common/log"
)

type Ticket interface {
	Grant() (tid string, err error)
	Access() string
}

type globalTicket struct {
	ticket   string
	provider TicketProvider
	mutex    sync.RWMutex
}

func (t *globalTicket) Grant() (tid string, err error) {
	tid, ticket, err := t.provider.Generate()
	if err != nil {
		return "", err
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.ticket = ticket

	return tid, nil
}

func (t *globalTicket) Access() string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.ticket
}

func initGlobalTicket(provider TicketProvider) Ticket {
	_, initTicket, err := provider.Generate()
	if err != nil {
		log.Logger.Fatal(err)
	}
	return &globalTicket{
		ticket:   initTicket,
		provider: provider,
		mutex:    sync.RWMutex{},
	}
}
