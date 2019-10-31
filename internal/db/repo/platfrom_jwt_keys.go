package repo

import (
	"context"

	"github.com/qilin/crm-api/internal/db/trx"

	"github.com/jinzhu/gorm"
	"github.com/qilin/crm-api/internal/db/domain"
)

type PlatformJWTKeyRepo struct {
	db *gorm.DB
}

func (a *PlatformJWTKeyRepo) All(ctx context.Context) ([]domain.PlatformJWTKeyItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		keys = []domain.PlatformJWTKeyItem{}
	)
	e := db.Find(&keys).Error
	return keys, e
}

func (a *PlatformJWTKeyRepo) Create(ctx context.Context, model *domain.PlatformJWTKeyItem) error {
	db := trx.Inject(ctx, a.db)
	return db.Save(model).Error
}

func (a *PlatformJWTKeyRepo) Get(ctx context.Context, alg, iss string) (*domain.PlatformJWTKeyItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		out = &domain.PlatformJWTKeyItem{}
		e   error
	)
	e = db.Where("alg=? AND iss=?", alg, iss).First(out).Error
	return out, e
}

func (a *PlatformJWTKeyRepo) Delete(ctx context.Context, item *domain.PlatformJWTKeyItem) error {
	db := trx.Inject(ctx, a.db)
	return db.Delete(item).Error
}

func NewPlatformJWTKeyRepo(db *gorm.DB) domain.PlatformJWTKeyRepo {
	return &PlatformJWTKeyRepo{db: db}
}
