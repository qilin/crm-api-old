package repo

import (
	"context"

	"github.com/qilin/crm-api/internal/db/trx"

	"github.com/jinzhu/gorm"
	"github.com/qilin/crm-api/internal/db/domain"
)

type ProductsRepo struct {
	db *gorm.DB
}

func (a *ProductsRepo) All(ctx context.Context, limit, offset int) ([]domain.ProductItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		keys = []domain.ProductItem{}
	)
	e := db.Find(&keys).Error
	return keys, e
}

func (a *ProductsRepo) Create(ctx context.Context, model *domain.ProductItem) error {
	db := trx.Inject(ctx, a.db)
	return db.Save(model).Error
}

func (a *ProductsRepo) Get(ctx context.Context, id string) (*domain.ProductItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		out = &domain.ProductItem{}
		e   error
	)
	e = db.Where("id=?", id).First(out).Error
	return out, e
}

func (a *ProductsRepo) Delete(ctx context.Context, item *domain.ProductItem) error {
	db := trx.Inject(ctx, a.db)
	return db.Delete(item).Error
}

func NewProductsRepo(db *gorm.DB) domain.ProductsRepo {
	return &ProductsRepo{db: db}
}
