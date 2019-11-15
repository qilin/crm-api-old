package utils

import (
	"crypto/ecdsa"
	"time"

	"github.com/pascaldekloe/jwt"
)

const (
	qilinProductUUIDKey = "qilinProductUUID"
	userIDKey           = "userID"
)

func IssueJWT(sub, iss, userId, qilinProductUUID string, exp int, pk *ecdsa.PrivateKey) ([]byte, error) {
	claims := jwt.Claims{
		Set: map[string]interface{}{},
	}
	claims.Subject = sub
	claims.Issuer = iss
	claims.Set[userIDKey] = userId
	claims.Set[qilinProductUUIDKey] = qilinProductUUID

	if exp > 0 {
		now := time.Now().Round(time.Second)
		claims.Expires = jwt.NewNumericTime(now.Add(time.Duration(exp) * time.Minute))
	}

	// issue a JWT
	return claims.ECDSASign(jwt.ES512, pk)
}
