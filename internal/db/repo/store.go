package repo

import (
	"context"

	"github.com/qilin/crm-api/internal/db/trx"

	"github.com/jinzhu/gorm"
	"github.com/qilin/crm-api/internal/db/domain"
)

type StoreRepo struct {
	db *gorm.DB
}

func (a *StoreRepo) Create(ctx context.Context, model *domain.StoreItem) error {
	db := trx.Inject(ctx, a.db)
	return db.Save(model).Error
}

func (a *StoreRepo) Get(ctx context.Context, id int) (*domain.StoreItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		out = &domain.StoreItem{}
		e   error
	)
	e = db.Where("id=?", id).First(out).Error
	return out, e
}

func NewPlatformRepo(db *gorm.DB) domain.StoreRepo {
	return &StoreRepo{db: db}
}
