package repo

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/qilin/crm-api/internal/db/domain"
)

type PlatformJWTKeyRepo struct {
	db *gorm.DB
}

func (p *PlatformJWTKeyRepo) All(ctx context.Context) ([]domain.PlatformJWTKeyItem, error) {
	return nil, nil
}

func (p *PlatformJWTKeyRepo) Create(ctx context.Context, model *domain.PlatformJWTKeyItem) error {
	return nil
}

func (p *PlatformJWTKeyRepo) Get(ctx context.Context, alg, iss string) (*domain.PlatformJWTKeyItem, error) {
	return &domain.PlatformJWTKeyItem{}, nil
}

func (p *PlatformJWTKeyRepo) Delete(ctx context.Context, item *domain.PlatformJWTKeyItem) error {
	return nil
}

func NewPlatformJWTKeyRepo(db *gorm.DB) domain.PlatformJWTKeyRepo {
	return &PlatformJWTKeyRepo{db: db}
}
