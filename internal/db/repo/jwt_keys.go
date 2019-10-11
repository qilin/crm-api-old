package repo

import (
	"context"

	"github.com/qilin/crm-api/internal/db/domain"
	"github.com/qilin/crm-api/internal/db/trx"

	"github.com/jinzhu/gorm"
)

type JWTKeysRepo struct {
	db *gorm.DB
}

// Get All Keys
func (a *JWTKeysRepo) All(ctx context.Context) ([]domain.JWTKeysItem, error) {
	var (
		out []domain.JWTKeysItem
		e   error
	)
	db := trx.Inject(ctx, a.db)
	e = db.Find(&out).Error
	return out, e
}

// Create
func (a *JWTKeysRepo) Create(ctx context.Context, model *domain.JWTKeysItem) error {
	db := trx.Inject(ctx, a.db)
	return db.Save(model).Error
}

// Get
func (a *JWTKeysRepo) Get(ctx context.Context, alg, iss string) (*domain.JWTKeysItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		out = &domain.JWTKeysItem{}
		e   error
	)
	e = db.Where("alg=? AND iss=?", alg, iss).First(out).Error
	return out, e
}

// Delete
func (a *JWTKeysRepo) Delete(ctx context.Context, item *domain.JWTKeysItem) error {
	db := trx.Inject(ctx, a.db)
	return db.Unscoped().Delete(item).Error
}

func NewJwtKeysRepo(db *gorm.DB) domain.JWTKeysRepo {
	return &JWTKeysRepo{db: db}
}
