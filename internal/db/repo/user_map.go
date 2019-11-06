package repo

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/qilin/crm-api/internal/db/domain"
	"github.com/qilin/crm-api/internal/db/trx"
)

type UserMapRepo struct {
	db *gorm.DB
}

func (a *UserMapRepo) Create(ctx context.Context, model *domain.UserMapItem) error {
	db := trx.Inject(ctx, a.db)
	return db.Save(model).Error
}

func (a *UserMapRepo) Get(ctx context.Context, id int) (*domain.UserMapItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		out = &domain.UserMapItem{}
		e   error
	)
	e = db.Where("id=?", id).First(out).Error
	return out, e
}

func (a *UserMapRepo) FindByExternalID(ctx context.Context, platformId int, externalID string) (*domain.UserMapItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		out = &domain.UserMapItem{}
		e   error
	)
	e = db.Where("platform_id = ? AND external_id=?", externalID).First(out).Error
	return out, e
}

func (a *UserMapRepo) IsExists(ctx context.Context, id string) (bool, error) {
	var count int
	db := trx.Inject(ctx, a.db)
	e := db.Model(&domain.UserMapItem{}).Where("id=?", id).Count(&count).Error
	if e != nil {
		return false, e
	}
	return count > 0, nil
}

func (a *UserMapRepo) IsExistsWithExternalID(ctx context.Context, external_id string) (bool, error) {
	var count int
	db := trx.Inject(ctx, a.db)
	e := db.Model(&domain.UserMapItem{}).Where("external_id=?", external_id).Count(&count).Error
	if e != nil {
		return false, e
	}
	return count > 0, nil
}

func NewUserMapRepo(db *gorm.DB) domain.UserMapRepo {
	return &UserMapRepo{db: db}
}
