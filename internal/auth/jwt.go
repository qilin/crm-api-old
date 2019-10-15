package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var privateKey = []byte(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIDgMLiKqAZ1l3sNvFTXph9opUsoijGqItCUP8K6lLp9KoAoGCCqGSM49
AwEHoUQDQgAEJdLA2Q1fMTdr6a1mg3oCjwxJnj6wi5eiweQa2c7oBqF0RIrsNzAe
pS47rxHaOQl7REVI1btuJ+fITxhGY5r2wA==
-----END EC PRIVATE KEY-----`)

var publicKey = []byte(`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEJdLA2Q1fMTdr6a1mg3oCjwxJnj6w
i5eiweQa2c7oBqF0RIrsNzAepS47rxHaOQl7REVI1btuJ+fITxhGY5r2wA==
-----END PUBLIC KEY-----`)

// GenerateJWT returns token signed by private key with filled fields
func GenerateJWT(userId string) (string, error) {

	type CustomClaims struct {
		Hasura map[string]interface{} `json:"https://hasura.io/jwt/claims"`
		jwt.StandardClaims
	}

	var now = time.Now()

	claims := &CustomClaims{
		Hasura: map[string]interface{}{
			"x-hasura-default-role":  "owner",
			"x-hasura-allowed-roles": []string{"owner", "admin", "merchant", "developer", "supporter"},
		},
		StandardClaims: jwt.StandardClaims{
			Issuer:    "qilin.protocol.one",
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(1 * time.Minute).Unix(),
			Subject:   userId,
		},
	}

	var (
		sPrivateKey    interface{}
		err            error
		signinigMethod jwt.SigningMethod
	)

	signinigMethod = jwt.SigningMethodES256
	sPrivateKey, err = jwt.ParseECPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(signinigMethod, claims)
	ss, err := token.SignedString(sPrivateKey)

	if err != nil {
		return "", err
	}

	return ss, nil
}

func ValidateJWT(token string) bool {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return jwt.ParseECPublicKeyFromPEM(publicKey)
	})

	if err == nil && t.Claims.Valid() == nil {
		return true
	}

	return false
}
