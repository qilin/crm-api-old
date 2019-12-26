package repo

import (
	"context"

	"github.com/qilin/crm-api/internal/db/trx"

	"github.com/jinzhu/gorm"
	"github.com/qilin/crm-api/internal/db/domain"
)

type Users struct {
	db *gorm.DB
}

func (u *Users) Create(ctx context.Context, model *domain.UsersItem) error {
	db := trx.Inject(ctx, u.db)
	return db.Save(model).Error
}

func (u *Users) Get(ctx context.Context, id int) (*domain.UsersItem, error) {
	db := trx.Inject(ctx, u.db)
	var (
		out = &domain.UsersItem{}
		e   error
	)
	e = db.Where("id=?", id).First(out).Error
	return out, e
}

func (u *Users) FindByEmail(ctx context.Context, email string) (*domain.UsersItem, error) {
	db := trx.Inject(ctx, u.db)
	var (
		out = &domain.UsersItem{}
		e   error
	)
	e = db.Where("email=?", email).First(out).Error
	return out, e
}

func (u *Users) FindByPhone(ctx context.Context, phone string) (*domain.UsersItem, error) {
	db := trx.Inject(ctx, u.db)
	var (
		out = &domain.UsersItem{}
		e   error
	)
	e = db.Where("phone=?", phone).First(out).Error
	return out, e
}

func (u *Users) IsExistsEmail(ctx context.Context, email string) (bool, error) {
	var count int
	db := trx.Inject(ctx, u.db)
	e := db.Model(&domain.UsersItem{}).Where("email=?", email).Count(&count).Error
	if e != nil {
		return false, e
	}
	return count > 0, nil
}

func (u *Users) IsExistsPhone(ctx context.Context, phone string) (bool, error) {
	var count int
	db := trx.Inject(ctx, u.db)
	e := db.Model(&domain.UsersItem{}).Where("phone=?", phone).Count(&count).Error
	if e != nil {
		return false, e
	}
	return count > 0, nil
}

func NewUsersRepo(db *gorm.DB) domain.UsersRepo {
	return &Users{db: db}
}
