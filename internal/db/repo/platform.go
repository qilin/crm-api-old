package repo

import (
	"context"

	"github.com/qilin/crm-api/internal/db/trx"

	"github.com/jinzhu/gorm"
	"github.com/qilin/crm-api/internal/db/domain"
)

type PlatformRepo struct {
	db *gorm.DB
}

func (a *PlatformRepo) Create(ctx context.Context, model *domain.PlatformItem) error {
	db := trx.Inject(ctx, a.db)
	return db.Save(model).Error
}

func (a *PlatformRepo) Get(ctx context.Context, id int) (*domain.PlatformItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		out = &domain.PlatformItem{}
		e   error
	)
	e = db.Where("id=?", id).First(out).Error
	return out, e
}

func NewPlatformRepo(db *gorm.DB) domain.PlatformRepo {
	return &PlatformRepo{db: db}
}
