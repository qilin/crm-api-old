package repo

import (
	"context"

	"github.com/qilin/crm-api/internal/db/domain"
	"github.com/qilin/crm-api/internal/db/trx"

	"github.com/jinzhu/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

// Create
func (a *UserRepo) Create(ctx context.Context, model *domain.UserItem) error {
	db := trx.Inject(ctx, a.db)
	return db.Save(model).Error
}

// FindByExternalID
func (a *UserRepo) FindByExternalID(ctx context.Context, externalID string) (*domain.UserItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		out = &domain.UserItem{}
		e   error
	)
	e = db.Where("external_id=?", externalID).First(out).Error
	return out, e
}

// FindByEmail
func (a *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.UserItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		out = &domain.UserItem{}
		e   error
	)
	e = db.Where("email=?", email).First(out).Error
	return out, e
}

func (a *UserRepo) Get(ctx context.Context, id int) (*domain.UserItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		out = &domain.UserItem{}
		e   error
	)
	e = db.Where("id=?", id).First(out).Error
	return out, e
}

func (a *UserRepo) IsExistsEmail(ctx context.Context, email string) (bool, error) {
	var count int
	db := trx.Inject(ctx, a.db)
	e := db.Model(&domain.UserItem{}).Where("email=?", email).Count(&count).Error
	if e != nil {
		return false, e
	}
	return count > 0, nil
}

func NewUserRepo(db *gorm.DB) domain.UserRepo {
	return &UserRepo{db: db}
}
