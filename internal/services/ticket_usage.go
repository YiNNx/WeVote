package services

// func ConsumeTicketUsage(ticketID string, times int) error {
// 	count, err := models.IncrTicketUsageCount(ticketID, times)
// 	if err != nil {
// 		return err
// 	}
// 	if count > config.C.Ticket.UpperLimit {
// 		models.DecrTicketUsageCount(ticketID, times)
// 		return errors.New("tmp")
// 	}
// 	return nil
// }

// func InitTicketUsage(ticketID string) error {
// 	return models.InitializeKeyUsageCount(ticketID, config.C.Ticket.Expiration.Duration)
// }
