package auth

import (
	"crypto"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/qilin/crm-api/internal/db/domain"
)

const AuthAudience = "auth"
const HasuraAPIAudience = "hasura"

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
	// required by hasura
	DefaultRole  string   `json:"x-hasura-default-role"`
	AllowedRoles []string `json:"x-hasura-allowed-roles"`
	// custom, all must be strings
	UserID string `json:"x-hasura-user-id,omitempty"`
}

type AccessTokenClaims struct {
	SessionClaims `json:"https://qilin.protocol.one/claims"`
	jwt.StandardClaims
}

func (c *AccessTokenClaims) Valid() error {
	if !c.VerifyAudience(HasuraAPIAudience, true) {
		return fmt.Errorf("invalid token audience")
	}
	return c.StandardClaims.Valid()
}

func NewAccessClaims(user *domain.UserItem) *AccessTokenClaims {
	var now = time.Now()
	return &AccessTokenClaims{
		SessionClaims: SessionClaims{
			DefaultRole:  user.Role,
			AllowedRoles: []string{user.Role},
			UserID:       strconv.Itoa(user.ID),
		},
		StandardClaims: jwt.StandardClaims{
			Issuer:    "https://qilin.protocol.one",
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(time.Hour).Unix(),
			Subject:   user.ExternalID,
			Audience:  HasuraAPIAudience, // not validated in hasura
		},
	}
}

type RefreshTokenClaims struct {
	jwt.StandardClaims
}

func (c *RefreshTokenClaims) Valid() error {
	if !c.VerifyAudience(AuthAudience, true) {
		return fmt.Errorf("invalid token audience")
	}
	return c.StandardClaims.Valid()
}

func NewRefreshClaims(user *domain.UserItem) *RefreshTokenClaims {
	var now = time.Now()
	return &RefreshTokenClaims{
		jwt.StandardClaims{
			Issuer:    "https://qilin.protocol.one",
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(refreshLifeTime).Unix(),
			Subject:   user.ExternalID,
			Audience:  AuthAudience, // used and validated
		},
	}
}
