package crypto

import (
	"crypto"

	jwt "github.com/dgrijalva/jwt-go"
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
