package sdk

import "github.com/qilin/crm-api/internal/db/domain"

type Repo struct {
	Platform       domain.PlatformRepo
	PlatformJWTKey domain.PlatformJWTKeyRepo
}
