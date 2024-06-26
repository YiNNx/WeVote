package ticket

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/YiNNx/WeVote/internal/config"
)

// the claims of ticket
type TicketClaims struct {
	IssuedAt  int64  `json:"iat"` // ticket issuance time
	ExpiresAt int64  `json:"exp"` // ticket expiration time, determined by ticket validity period
	SubjectID string `json:"sub"` // ticket id, generated by uuid, unique identifier for each ticket
}

// Valid implements the jwt.Claims interface for TicketClaims
func (c *TicketClaims) Valid() error {
	vErr := new(jwt.ValidationError)
	now := time.Now().Unix()

	// Check whether the ticket has expired
	if now > c.ExpiresAt {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("token is expired by %v", delta)
		vErr.Errors = jwt.ValidationErrorExpired
	}

	// Verify whether the issuance time of the ticket is valid
	if now < c.IssuedAt {
		vErr.Inner = fmt.Errorf("Token used before issued")
		vErr.Errors = jwt.ValidationErrorIssuedAt
	}

	if vErr.Errors == 0 {
		return nil
	}
	return vErr
}

// build and issue the ticket
func Generate() (ticketID string, ticketStr string, err error) {
	ticketID = uuid.New().String()
	ticket := jwt.Token{
		Method: jwt.SigningMethodHS256,
		Header: map[string]interface{}{
			"typ": "wevote-ticket",
			"alg": jwt.SigningMethodHS256.Alg(),
		},
		Claims: &TicketClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(config.C.Ticket.TTL.Duration).Unix(),
			SubjectID: ticketID,
		},
	}
	ticketStr, err = ticket.SignedString([]byte(config.C.Ticket.Secret))
	return ticketID, ticketStr, err
}

// ParseAndVerify parses the claims of the ticket and verifies the signature
func ParseAndVerify(ticketStr string) (*TicketClaims, error) {
	token, err := jwt.ParseWithClaims(ticketStr, &TicketClaims{}, func(t *jwt.Token) (interface{}, error) { return []byte(config.C.Ticket.Secret), nil })
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*TicketClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
