package repo

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/qilin/crm-api/internal/db/domain"
)

type PlatformRepo struct {
	db *gorm.DB
}

func (p *PlatformRepo) Create(ctx context.Context, model *domain.PlatformItem) error {
	return nil
}

func (p *PlatformRepo) Get(ctx context.Context, id int) (*domain.PlatformItem, error) {
	return &domain.PlatformItem{}, nil
}

func NewPlatformRepo(db *gorm.DB) domain.PlatformRepo {
	return &PlatformRepo{db: db}
}
