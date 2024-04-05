package services

import (
	"context"
	"sync"

	"github.com/YiNNx/WeVote/internal/common/errors"
	"github.com/YiNNx/WeVote/internal/common/log"
	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/models"
	"github.com/YiNNx/WeVote/internal/pkg/cron"
	"github.com/YiNNx/WeVote/internal/pkg/ticket"
)

func GetTicket() string {
	return gTicket.access()
}

func ParseAndVerifyTicket(ticketStr string) (ticketID string, err error) {
	claims, err := ticket.ParseAndVerify(ticketStr)
	if err != nil {
		return "", errors.ErrInvalidTicket.WithErrDetail(err)
	}
	ticketID = claims.SubjectID
	return ticketID, nil
}

func TicketConsume(ctx context.Context, tid string, count int) error {
	return gTicket.consume(ctx, tid, count)
}

// gTicket is the global Shared Ticket instance
var gTicket *sharedTicket

type sharedTicket struct {
	ticket string
	// Use read-write locks to deal with multi-threaded data race
	mutex sync.RWMutex
	// to count the limit
	models.Counter
}

// access externally exposed access method
func (t *sharedTicket) access() string {
	// Enable read lock
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.ticket
}

func (t *sharedTicket) consume(ctx context.Context, tid string, count int) error {
	return t.Counter.IncreaseBy(ctx, tid, count)
}

// issues a ticket and set the counter ttl
func (t *sharedTicket) grant() error {
	tid, err := t.generate()
	if err != nil {
		return err
	}
	err = t.Counter.Set(context.Background(), tid)
	if err != nil {
		return err
	}
	return nil
}

// issues a ticket and uses it to update the global ticket
func (t *sharedTicket) generate() (tid string, err error) {
	tid, ticket, err := ticket.Generate()
	if err != nil {
		return "", err
	}

	// Enable write lock
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.ticket = ticket

	return tid, nil
}

// init the global Ticket instance
func initGlobalSharedTicket() {
	gTicket = &sharedTicket{
		mutex: sync.RWMutex{},
		Counter: *models.NewTicketUsageCounter(
			config.C.Ticket.TTL.Duration,
			config.C.Ticket.UsageLimit,
		),
	}

	cronJobTicketGrant := cron.NewCronJob(
		config.C.Ticket.TTL.Duration,
		func() {
			err := gTicket.grant()
			if err != nil {
				log.Logger.Fatal(err)
			}
		})
	cronJobTicketGrant.Start()
}
