package ticket

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const tokenTypeTicket = "ticket"

type Ticket struct {
	jwt.Token
}

type TicketClaims struct {
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	SubjectId string `json:"sub"`
}

func (c *TicketClaims) Valid() error {
	vErr := new(jwt.ValidationError)
	now := time.Now().Unix()

	if now > c.ExpiresAt {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("token is expired by %v", delta)
		vErr.Errors = jwt.ValidationErrorExpired
	}

	if now < c.IssuedAt {
		vErr.Inner = fmt.Errorf("Token used before issued")
		vErr.Errors = jwt.ValidationErrorIssuedAt
	}

	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}

var Generator *generator

type generator struct {
	secret string
	spec   int
}

func (g *generator) Generate() (string, error) {
	ticket := Ticket{
		jwt.Token{
			Method: jwt.SigningMethodHS256,
			Header: map[string]interface{}{
				"typ": tokenTypeTicket,
				"alg": jwt.SigningMethodHS256.Alg(),
			},
			Claims: &TicketClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Duration(g.spec) * time.Second).Unix(),
				SubjectId: uuid.New().String(),
			},
		},
	}
	return ticket.SignedString(g.secret)
}

func InitGenerator(secret string, spec int) {
	Generator.secret = secret
	Generator.spec = spec
}
