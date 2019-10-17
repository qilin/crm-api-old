package auth

import (
	"crypto"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/qilin/crm-api/internal/db/domain"
)

// Keys =======================================================================

type KeyPair struct {
	Public  crypto.PublicKey
	Private crypto.PrivateKey
}

func NewKeyPair(public crypto.PublicKey, private crypto.PrivateKey) KeyPair {
	return KeyPair{public, private}
}

func NewKeyPairFromPEM(public, private string) (KeyPair, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(private))
	if err != nil {
		return KeyPair{}, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(public))
	if err != nil {
		return KeyPair{}, err
	}
	return NewKeyPair(publicKey, privateKey), nil
}

// Claims =====================================================================

type SessionClaims struct {
	// required by hasura
	DefaultRole  string   `json:"x-hasura-default-role"`
	AllowedRoles []string `json:"x-hasura-allowed-roles"`
	// custom, all must be strings
	UserID string `json:"x-hasura-user-id,omitempty"`
}

type TokenClaims struct {
	SessionClaims `json:"https://qilin.protocol.one/claims"`
	jwt.StandardClaims
}

func NewClaims(user *domain.UserItem) *TokenClaims {
	var now = time.Now()
	return &TokenClaims{
		SessionClaims: SessionClaims{
			DefaultRole:  "owner",
			AllowedRoles: []string{"owner", "admin", "merchant", "developer", "supporter"},
			UserID:       strconv.Itoa(user.ID),
		},
		StandardClaims: jwt.StandardClaims{
			Issuer:    "https://qilin.protocol.one",
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(time.Hour).Unix(),
			Subject:   user.ExternalID,
			Audience:  "some",
		},
	}
}
