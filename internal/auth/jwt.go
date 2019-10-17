package auth

import (
	"crypto"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/qilin/crm-api/internal/db/domain"
)

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

type TokenClaims struct {
	Hasura jwt.MapClaims `json:"https://qilin.protocol.one/claims"`
	jwt.StandardClaims
}

func NewClaims(user *domain.UserItem) *TokenClaims {
	var now = time.Now()
	return &TokenClaims{
		Hasura: map[string]interface{}{
			"x-hasura-default-role":  "owner",
			"x-hasura-allowed-roles": []string{"owner", "admin", "merchant", "developer", "supporter"},
			"x-hasura-user-id":       strconv.Itoa(user.ID),
		},
		StandardClaims: jwt.StandardClaims{
			Issuer:    "https://qilin.protocol.one",
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(30 * time.Minute).Unix(),
			Subject:   user.ExternalID,
			Audience:  "some",
		},
	}
}
