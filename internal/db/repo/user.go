package repo

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	domain2 "github.com/qilin/crm-api/internal/db/domain"
	trx2 "github.com/qilin/crm-api/internal/db/trx"

	"github.com/jinzhu/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

// Create
func (a *UserRepo) Create(ctx context.Context, model *domain2.UserItem) error {
	db := trx2.Inject(ctx, a.db)
	pwd, e := hashPassword(model.Password)
	if e != nil {
		return e
	}
	model.Password = pwd
	return db.Save(model).Error
}

// List
func (a *UserRepo) Get(ctx context.Context, email string, password string) (*domain2.UserItem, error) {
	db := trx2.Inject(ctx, a.db)
	var (
		out *domain2.UserItem
		e   error
	)
	e = db.Where("email=? AND password=?", email, password).First(out).Error
	return out, e
}

func (a *UserRepo) IsExistsEmail(ctx context.Context, email string) (bool, error) {
	var count int
	db := trx2.Inject(ctx, a.db)
	e := db.Model(&domain2.UserItem{}).Where("email=?", email).Count(&count).Error
	if e != nil {
		return false, e
	}
	return count > 0, nil
}

func NewUserRepo(db *gorm.DB) domain2.UserRepo {
	return &UserRepo{db: db}
}

func hashPassword(password string) (string, error) {
	hash, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), e
}
