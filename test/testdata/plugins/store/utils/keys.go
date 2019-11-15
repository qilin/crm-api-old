package utils

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/dgrijalva/jwt-go"
)

type KeyPair struct {
	Public  *ecdsa.PublicKey
	Private *ecdsa.PrivateKey
}

type RSAKeyPair struct {
	Public  *rsa.PublicKey
	Private *rsa.PrivateKey
}

func DecodePemECDSA(pemPriv, pemPub string) (KeyPair, error) {
	block, _ := pem.Decode([]byte(pemPriv))
	x509Encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		return KeyPair{}, err
	}

	blockPub, _ := pem.Decode([]byte(pemPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return KeyPair{}, err
	}
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return KeyPair{
		Public:  publicKey,
		Private: privateKey,
	}, nil
}

func DecodePemRSA(private, public string) (RSAKeyPair, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(private))
	if err != nil {
		return RSAKeyPair{}, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(public))
	if err != nil {
		return RSAKeyPair{}, err
	}

	return RSAKeyPair{
		Public:  publicKey,
		Private: privateKey,
	}, nil
}
