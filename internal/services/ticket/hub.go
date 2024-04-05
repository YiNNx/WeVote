package ticket

import (
	"sync"

	"github.com/YiNNx/WeVote/internal/common/log"
	"github.com/YiNNx/WeVote/internal/pkg/ticket"
)

// gTicket is the global Ticket instance
var gTicket *globalTicket

type globalTicket struct {
	ticket string
	// Use read-write locks to deal with multi-threaded data race
	mutex sync.RWMutex
}

// Access externally exposed access method
func Access() string {
	// Enable read lock
	gTicket.mutex.RLock()
	defer gTicket.mutex.RUnlock()

	return gTicket.ticket
}

// globalGrant issues a ticket and uses it to update the global ticket
func globalGrant() (tid string, err error) {
	tid, ticket, err := ticket.Generate()
	if err != nil {
		return "", err
	}

	// Enable write lock
	gTicket.mutex.Lock()
	defer gTicket.mutex.Unlock()

	gTicket.ticket = ticket

	return tid, nil
}

// init the global Ticket instance
func init() {
	_, initTicket, err := ticket.Generate()
	if err != nil {
		log.Logger.Fatal(err)
	}
	gTicket = &globalTicket{
		ticket: initTicket,
		mutex:  sync.RWMutex{},
	}
}
