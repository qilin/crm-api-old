package jwt

import (
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/qilin/crm-api/internal/db/domain"
)

// ProviderJwtVerifier
func ProviderJwtVerifier(repo domain.JWTKeysRepo, set provider.AwareSet) *JWTVerefier {
	return NewJWTVerifier(repo, set)
}
