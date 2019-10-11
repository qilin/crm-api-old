package jwt

import "github.com/qilin/crm-api/internal/db/domain"

// ProviderJwtVerifier
func ProviderJwtVerifier(repo domain.JWTKeysRepo) *JWTVerefier {
	return NewJWTVerifier(repo)
}
