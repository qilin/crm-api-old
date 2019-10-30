package auth

import (
	"crypto"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/qilin/crm-api/internal/db/domain"
)

const Audience = "qilin"

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

func (keys *KeyPair) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	return token.SignedString(keys.Private)
}

func (keys *KeyPair) Parse(rawToken string, claims jwt.Claims) error {
	_, err := jwt.ParseWithClaims(rawToken, claims, func(*jwt.Token) (interface{}, error) {
		return keys.Public, nil
	})
	return err
}

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
