package ticket

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/YiNNx/WeVote/internal/config"
)

type Ticket struct {
	jwt.Token
}

type Claims struct {
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	SubjectId string `json:"sub"`
}

func (c *Claims) Valid() error {
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

type Provider interface {
	Generate() (ticketID string, ticketStr string, err error)
}

type provider struct {
	secret     string
	expiration time.Duration
}

func (p *provider) Generate() (ticketID string, ticketStr string, err error) {
	ticketID = uuid.New().String()
	ticket := Ticket{
		jwt.Token{
			Method: jwt.SigningMethodHS256,
			Header: map[string]interface{}{
				"typ": "wevote-ticket",
				"alg": jwt.SigningMethodHS256.Alg(),
			},
			Claims: &Claims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(p.expiration).Unix(),
				SubjectId: ticketID,
			},
		},
	}
	ticketStr, err = ticket.SignedString(p.secret)
	return ticketID, ticketStr, err
}

func NewProvider(secret string, expiration time.Duration) Provider {
	return &provider{
		secret:     secret,
		expiration: expiration,
	}
}

func ParseAndVerifyTicket(ticketStr string) (ticketClaims *Claims, err error) {
	token, err := jwt.ParseWithClaims(ticketStr, &Claims{}, func(t *jwt.Token) (interface{}, error) { return config.C.Ticket.Secret, nil })
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
