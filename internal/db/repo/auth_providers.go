package repo

import (
	"context"

	"github.com/qilin/crm-api/internal/db/trx"

	"github.com/jinzhu/gorm"
	"github.com/qilin/crm-api/internal/db/domain"
)

type AuthProvider struct {
	db *gorm.DB
}

func (a *AuthProvider) Create(ctx context.Context, model *domain.AuthProviderItem) error {
	db := trx.Inject(ctx, a.db)
	return db.Save(model).Error
}

func (a *AuthProvider) Delete(ctx context.Context, model *domain.AuthProviderItem) error {
	db := trx.Inject(ctx, a.db)
	return db.Delete(model).Error
}

func (a *AuthProvider) Get(ctx context.Context, user_id int, provider, provider_id string) (*domain.AuthProviderItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		out = &domain.AuthProviderItem{}
		e   error
	)
	e = db.Where("user_id=? AND provider=? AND provider_id=?", user_id, provider, provider_id).First(out).Error
	return out, e
}

func (a *AuthProvider) GetByUserId(ctx context.Context, user_id int) ([]domain.AuthProviderItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		items = []domain.AuthProviderItem{}
	)
	e := db.Find(&items).Error
	return items, e
}

func NewAuthProvider(db *gorm.DB) domain.AuthProviderRepo {
	return &AuthProvider{db: db}
}
