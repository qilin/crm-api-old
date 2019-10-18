package jwt

import (
	"strings"

	"github.com/ProtocolONE/go-core/v2/pkg/provider"

	"golang.org/x/net/context"

	"github.com/pascaldekloe/jwt"
	"github.com/qilin/crm-api/internal/db/domain"
)

const (
	Prefix = "internal.jwt"
)

type JWTVerefier struct {
	repo        domain.JWTKeysRepo
	keys        []*domain.JWTKeysItem
	keyRegister jwt.KeyRegister
}

func (j *JWTVerefier) Check(token string) (*jwt.Claims, error) {
	return j.keyRegister.Check([]byte(token))
}

func (j *JWTVerefier) LoadKeys() error {
	keys, err := j.repo.All(context.Background())
	if err != nil {
		return err
	}
	for i := range keys {
		switch strings.ToLower(keys[i].KeyType) {
		case "pem":
			j.keyRegister.LoadPEM([]byte(keys[i].Key), []byte(""))
		case "jwk":
			j.keyRegister.LoadJWK([]byte(keys[i].Key))
		}
	}
	return nil
}

func NewJWTVerifier(repo domain.JWTKeysRepo, set provider.AwareSet) *JWTVerefier {
	j := &JWTVerefier{
		repo: repo,
	}
	if j.LoadKeys() != nil {
		set.L().Emergency("can't load jwt keys")
	}
	return j
}
