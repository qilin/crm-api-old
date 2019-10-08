package repo

import (
	"context"

	domain2 "github.com/qilin/crm-api/internal/db/domain"
	trx2 "github.com/qilin/crm-api/internal/db/trx"

	"github.com/jinzhu/gorm"
)

type JWTKeysRepo struct {
	db *gorm.DB
}

// Get All Keys
func (a *JWTKeysRepo) All(ctx context.Context) ([]domain2.JWTKeysItem, error) {
	var (
		out []domain2.JWTKeysItem
		e   error
	)
	db := trx2.Inject(ctx, a.db)
	e = db.Find(out).Error
	return out, e
}

// Create
func (a *JWTKeysRepo) Create(ctx context.Context, model *domain2.JWTKeysItem) error {
	db := trx2.Inject(ctx, a.db)
	return db.Save(model).Error
}

// Get
func (a *JWTKeysRepo) Get(ctx context.Context, alg, iss string) (*domain2.JWTKeysItem, error) {
	db := trx2.Inject(ctx, a.db)
	var (
		out = &domain2.JWTKeysItem{}
		e   error
	)
	e = db.Where("alg=? AND iss=?", alg, iss).First(out).Error
	return out, e
}

// Delete
func (a *JWTKeysRepo) Delete(ctx context.Context, item *domain2.JWTKeysItem) error {
	db := trx2.Inject(ctx, a.db)
	return db.Unscoped().Delete(item).Error
}

func NewJwtKeysRepo(db *gorm.DB) domain2.JWTKeysRepo {
	return &JWTKeysRepo{db: db}
}
