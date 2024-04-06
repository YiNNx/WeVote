package services

import (
	"context"
	"time"

	"github.com/YiNNx/WeVote/internal/common/errors"
	"github.com/YiNNx/WeVote/internal/common/log"
	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/models"
	"github.com/YiNNx/WeVote/internal/pkg/ticket"
	"github.com/YiNNx/WeVote/pkg/cronjob"
)

func GetTicket(ctx context.Context) (string, error) {
	return srvTicket.access(ctx)
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
	return srvTicket.consume(ctx, tid, count)
}

var srvTicket *serviceTicket

type serviceTicket struct {
	ticket  models.GlobalTicket
	counter models.Counter
}

func (t *serviceTicket) access(ctx context.Context) (string, error) {
	return t.ticket.Access(ctx)
}

func (t *serviceTicket) consume(ctx context.Context, tid string, count int) error {
	return t.counter.IncreaseBy(ctx, tid, count)
}

func (t *serviceTicket) refresh(ctx context.Context) error {
	tid, ticket, err := ticket.Generate()
	if err != nil {
		return err
	}

	err = t.counter.Set(ctx, tid)
	if err != nil {
		return err
	}
	err = t.ticket.Save(ctx, ticket)
	if err != nil {
		return err
	}
	log.Logger.Info("refreshed the ticket: " + ticket)
	return nil
}

func initServiceTicket() {
	srvTicket = &serviceTicket{
		ticket: models.InitGlobalTicket(
			config.C.Ticket.TTL.Duration,
		),
		counter: models.InitTicketUsageCounter(
			config.C.Ticket.TTL.Duration+time.Duration(1)*time.Second,
			config.C.Ticket.UsageLimit,
		),
	}
}

func StartCronJobTicketRefresh() {
	srvTicket = &serviceTicket{
		ticket: models.InitGlobalTicket(
			config.C.Ticket.TTL.Duration,
		),
		counter: models.InitTicketUsageCounter(
			config.C.Ticket.TTL.Duration+time.Duration(1)*time.Second,
			config.C.Ticket.UsageLimit,
		),
	}
	cronJobTicketRefresh := cronjob.NewCronJob(
		config.C.Ticket.TTL.Duration,
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()
			err := srvTicket.refresh(ctx)
			if err != nil {
				log.Logger.Error(err)
			}
		})
	cronJobTicketRefresh.Func()
	cronJobTicketRefresh.Start()
}
