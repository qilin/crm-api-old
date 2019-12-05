package repo

import (
	"context"
	"fmt"

	"github.com/qilin/crm-api/internal/db/trx"

	"github.com/jinzhu/gorm"
	"github.com/qilin/crm-api/internal/db/domain"
	"github.com/qilin/crm-api/internal/db/domain/store"
)

type StoreJWTKeyRepo struct {
	db *gorm.DB
}

func (a *StoreJWTKeyRepo) All(ctx context.Context) ([]domain.StoreJWTKeyItem, error) {
	fmt.Println("touch games")
	r := NewGamesRepo(a.db)
	err := r.Insert(ctx, &store.Game{
		ID: "fa14b399-ae9b-4111-9c7f-0f1fe2cc1eb6",
	})
	if err != nil {
		fmt.Println("error:", err)
	}
	err = r.Delete(ctx,
		"fa14b399-ae9b-4111-9c7f-0f1fe2cc1eb6",
	)
	if err != nil {
		fmt.Println("error:", err)
	}

	db := trx.Inject(ctx, a.db)
	var (
		keys = []domain.StoreJWTKeyItem{}
	)
	e := db.Find(&keys).Error
	return keys, e
}

func (a *StoreJWTKeyRepo) Create(ctx context.Context, model *domain.StoreJWTKeyItem) error {
	db := trx.Inject(ctx, a.db)
	return db.Save(model).Error
}

func (a *StoreJWTKeyRepo) Get(ctx context.Context, alg, iss string) (*domain.StoreJWTKeyItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		out = &domain.StoreJWTKeyItem{}
		e   error
	)
	e = db.Where("alg=? AND iss=?", alg, iss).First(out).Error
	return out, e
}

func (a *StoreJWTKeyRepo) GetByKID(ctx context.Context, kid string) (*domain.StoreJWTKeyItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		out = &domain.StoreJWTKeyItem{}
		e   error
	)
	e = db.Where("kid=?", kid).First(out).Error
	return out, e
}

func (a *StoreJWTKeyRepo) GetByIss(ctx context.Context, iss string) (*domain.StoreJWTKeyItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		out = &domain.StoreJWTKeyItem{}
		e   error
	)
	e = db.Where("iss=?", iss).First(out).Error
	return out, e
}

func (a *StoreJWTKeyRepo) Delete(ctx context.Context, item *domain.StoreJWTKeyItem) error {
	db := trx.Inject(ctx, a.db)
	return db.Delete(item).Error
}

func NewPlatformJWTKeyRepo(db *gorm.DB) domain.StoreJWTKeyRepo {
	return &StoreJWTKeyRepo{db: db}
}
