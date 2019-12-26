package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"

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

type ECKeyPair struct {
	Public  crypto.PublicKey
	Private crypto.PrivateKey
}

func NewECKeyPair(public crypto.PublicKey, private crypto.PrivateKey) ECKeyPair {
	return ECKeyPair{public, private}
}

func NewECKeyPairFromPEM(public, private string) (ECKeyPair, error) {
	privateKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(private))
	if err != nil {
		return ECKeyPair{}, err
	}
	publicKey, err := jwt.ParseECPublicKeyFromPEM([]byte(public))
	if err != nil {
		return ECKeyPair{}, err
	}
	return NewECKeyPair(publicKey, privateKey), nil
}

func (keys *ECKeyPair) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(keys.Private)
}

func (keys *ECKeyPair) Parse(rawToken string, claims jwt.Claims) error {
	_, err := jwt.ParseWithClaims(rawToken, claims, func(*jwt.Token) (interface{}, error) {
		return keys.Public, nil
	})
	return err
}

type ECDSAKeyPair struct {
	Public  *ecdsa.PublicKey
	Private *ecdsa.PrivateKey
}

func NewECDSAKeyPairFromPEM(public, private string) (ECDSAKeyPair, error) {
	block, _ := pem.Decode([]byte(private))
	x509Encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		return ECDSAKeyPair{}, err
	}

	blockPub, _ := pem.Decode([]byte(public))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return ECDSAKeyPair{}, err
	}
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return NewECDSAKeyPair(publicKey, privateKey)
}

func NewECDSAKeyPair(public *ecdsa.PublicKey, private *ecdsa.PrivateKey) (ECDSAKeyPair, error) {
	return ECDSAKeyPair{
		Public:  public,
		Private: private,
	}, nil
}

func (keys *ECDSAKeyPair) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES384, claims)
	return token.SignedString(keys.Private)
}

func (keys *ECDSAKeyPair) Parse(rawToken string, claims jwt.Claims) error {
	_, err := jwt.ParseWithClaims(rawToken, claims, func(*jwt.Token) (interface{}, error) {
		return keys.Public, nil
	})
	return err
}
