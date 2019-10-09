package repo

import (
	"context"

	"golang.org/x/crypto/bcrypt"

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
	pwd, e := hashPassword(model.Password)
	if e != nil {
		return e
	}
	model.Password = pwd
	return db.Save(model).Error
}

// List
func (a *UserRepo) Get(ctx context.Context, email string, password string) (*domain.UserItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		out *domain.UserItem
		e   error
	)
	pwd, e := hashPassword(password)
	if e != nil {
		return out, e
	}
	e = db.Where("email=? AND password=?", email, pwd).First(out).Error
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

func hashPassword(password string) (string, error) {
	hash, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), e
}
