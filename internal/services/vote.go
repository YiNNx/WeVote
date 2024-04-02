package services

import (
	"context"
	"errors"

	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/models"
)

type UserSet map[string]struct{}

func Vote(ticketID string, users UserSet) error {
	rtx := models.BeginRedisTx()
	userCount := len(users)
	count, err := rtx.IncrTicketUsageCount(ticketID, userCount)
	if err != nil {
		rtx.Discard()
		return err
	}
	if count > config.C.Ticket.UpperLimit {
		rtx.Discard()
		return errors.New("tmp")
	}

	for user := range users {
		_, err = rtx.IncrVoteCount(user)
		if err != nil {
			rtx.Discard()
			return err
		}
	}
	_, err = rtx.Exec(context.Background())
	return err
}

// func VoteUserBatch(usernames []string) error {

// 	for _, username := range usernames {
// 		err := vote(username)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func vote(username string) error {
// 	_, err := models.IncrVoteCount(username)
// 	return err
// }
