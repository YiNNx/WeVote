package services

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/YiNNx/WeVote/internal/common/errors"
	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/models"
)

type UserSet map[string]struct{}

func Vote(ctx context.Context, ticketID string, users UserSet) error {
	rtx := models.BeginRedisTx(ctx)
	userCount := len(users)
	count, err := rtx.IncrTicketUsageCount(ticketID, userCount)
	if err != nil {
		rtx.Discard()
		return errors.DataUpdateFailed.WithErrDetail(err)
	}
	if count > config.C.Ticket.UpperLimit {
		rtx.Discard()
		return errors.TicketUsageLimitExceed
	}

	userList := make([]string, 0, len(users))
	for user := range users {
		_, err = rtx.IncrVoteCount(user)
		if err != nil {
			rtx.Discard()
			return errors.DataUpdateFailed.WithErrDetail(err)
		}
		userList = append(userList, user)
	}

	err = rtx.RecordVoteCountModified(userList)
	if err != nil {
		rtx.Discard()
		return errors.DataUpdateFailed.WithErrDetail(err)
	}

	_, err = rtx.Exec(context.Background())
	if err != nil {
		return errors.DataUpdateFailed.WithErrDetail(err)
	}
	return nil
}

func GetVoteCount(ctx context.Context, user string) (int, error) {
	count, err := models.GetVoteCount(ctx, user)
	if err == redis.Nil {
		return 0, errors.UserNotFound // TODO:
	}
	if err != nil {
		return 0, errors.DataLoadFailed.WithErrDetail(err)
	}
	res, err := strconv.Atoi(count)
	if err != nil {
		return 0, errors.DataLoadFailed.WithErrDetail(err)
	}
	return res, nil
}
