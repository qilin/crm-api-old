package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/qilin/crm-api/internal/db/domain"
)

const Audience = "qilin"

// Claims =====================================================================

type SessionClaims struct {
	UserID   string `json:"user-id,omitempty"`
	TenantID string `json:"tenant-id,omitempty"`
	Role     string `json:"role,omitempty"`
}

type AccessTokenClaims struct {
	SessionClaims `json:"https://qilin.protocol.one/claims"`
	jwt.StandardClaims
}

func (c *AccessTokenClaims) Valid() error {
	if !c.VerifyAudience(Audience, true) {
		return fmt.Errorf("invalid token audience")
	}
	return c.StandardClaims.Valid()
}

func NewAccessClaims(user *domain.UserItem) *AccessTokenClaims {
	var now = time.Now()
	return &AccessTokenClaims{
		SessionClaims: SessionClaims{
			UserID:   strconv.Itoa(user.ID),
			TenantID: strconv.Itoa(user.TenantID),
			Role:     user.Role,
		},
		StandardClaims: jwt.StandardClaims{
			Issuer:   "https://qilin.protocol.one",
			IssuedAt: now.Unix(),
			Subject:  user.ExternalID,
			Audience: Audience, // not validated in hasura
		},
	}
}
