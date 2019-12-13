package authentication

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/qilin/crm-api/internal/db/domain"
)

type SessionClaims struct {
	User User
	jwt.StandardClaims
}

// todo: Audience
func (c *SessionClaims) Valid() error {
	if !c.VerifyAudience("Audience", true) {
		return fmt.Errorf("invalid token audience")
	}
	return c.StandardClaims.Valid()
}

func NewAccessClaims(user *domain.UsersItem) *SessionClaims {
	var now = time.Now()
	return &SessionClaims{
		User: User{
			ID:       user.ID,
			Language: user.Language,
		},
		StandardClaims: jwt.StandardClaims{
			Issuer:   "https://qilin.protocol.one", // todo: to config
			IssuedAt: now.Unix(),
			Subject:  "",         // todo: to config
			Audience: "Audience", // todo: to config
		},
	}
}
