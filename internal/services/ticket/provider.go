package ticket

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Token struct {
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

type TicketProvider interface {
	Generate() (ticketID string, ticketStr string, err error)
}

type provider struct {
	secret     string
	expiration time.Duration
}

func (p *provider) Generate() (ticketID string, ticketStr string, err error) {
	ticketID = uuid.New().String()
	ticket := Token{
		jwt.Token{
			Method: jwt.SigningMethodHS256,
			Header: map[string]interface{}{
				"typ": "wevote-ticket",
				"alg": jwt.SigningMethodHS256.Alg(),
			},
			Claims: &TicketClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(p.expiration).Unix(),
				SubjectId: ticketID,
			},
		},
	}
	ticketStr, err = ticket.SignedString(p.secret)
	return ticketID, ticketStr, err
}

func newTicketProvider(secret string, expiration time.Duration) TicketProvider {
	return &provider{
		secret:     secret,
		expiration: expiration,
	}
}

type TicketParser interface {
	ParseAndVerify(ticketStr string) (ticketClaims *TicketClaims, err error)
}

type parser struct {
	secret string
}

func (p *parser) ParseAndVerify(ticketStr string) (ticketClaims *TicketClaims, err error) {
	token, err := jwt.ParseWithClaims(ticketStr, &TicketClaims{}, func(t *jwt.Token) (interface{}, error) { return p.secret, nil })
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*TicketClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

func newParser(secret string) TicketParser {
	return &parser{
		secret: secret,
	}
}
