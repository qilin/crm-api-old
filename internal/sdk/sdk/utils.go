package sdk

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
)

type KeyPair struct {
	Public  *ecdsa.PublicKey
	Private *ecdsa.PrivateKey
}

func decodePemECDSA(pemPriv, pemPub string) (KeyPair, error) {
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
