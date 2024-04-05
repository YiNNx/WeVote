package vote

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/YiNNx/WeVote/internal/common/errors"
	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/models"
)

func ProcessVote(ctx context.Context, ticketID string, users UsernameSet) error {
	rtx := models.BeginRedisTx(ctx)
	userCount := len(users)
	count, err := rtx.IncrTicketUsageCount(ticketID, userCount)
	if err != nil {
		rtx.Discard()
		return errors.ErrDataUpdate.WithErrDetail(err)
	}
	if count > config.C.Ticket.UpperLimit {
		rtx.Discard()
		return errors.ErrTicketUsageExceed
	}

	userList := make([]string, 0, len(users))
	for user := range users {
		_, err = rtx.IncrVoteCount(user)
		if err != nil {
			rtx.Discard()
			return errors.ErrDataUpdate.WithErrDetail(err)
		}
		userList = append(userList, user)
	}

	err = rtx.RecordUserModified(userList)
	if err != nil {
		rtx.Discard()
		return errors.ErrDataUpdate.WithErrDetail(err)
	}

	_, err = rtx.Exec(context.Background())
	if err != nil {
		return errors.ErrDataUpdate.WithErrDetail(err)
	}
	return nil
}

func GetVoteCount(ctx context.Context, user string) (int, error) {
	count, err := models.GetVoteCountByCache(ctx, user)
	if err == redis.Nil {
		return 0, errors.ErrInvalidUsernameExisted
	}
	if err != nil {
		return 0, errors.ErrDataLoad.WithErrDetail(err)
	}
	res, err := strconv.Atoi(count)
	if err != nil {
		return 0, errors.ErrDataLoad.WithErrDetail(err)
	}
	return res, nil
}
