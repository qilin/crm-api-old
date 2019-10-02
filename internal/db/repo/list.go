package repo

import (
	"context"
	"strings"

	domain2 "github.com/qilin/crm-api/internal/db/domain"
	trx2 "github.com/qilin/crm-api/internal/db/trx"

	"github.com/jinzhu/gorm"
)

type ListRepo struct {
	db *gorm.DB
}

// Create
func (a *ListRepo) Create(ctx context.Context, model *domain2.ListItem) error {
	db := trx2.Inject(ctx, a.db)
	return db.Save(model).Error
}

// List
func (a *ListRepo) List(ctx context.Context, projection []string, cursor *domain2.Cursor, order domain2.Order, search string) ([]*domain2.ListItem, error) {
	//
	cursor.Init()
	db := trx2.Inject(ctx, a.db)
	//
	var (
		out []*domain2.ListItem
		e   error
	)
	if len(projection) > 0 {
		db = db.Select(projection)
	}
	db = db.Where("name=?", search)
	e = db.Model(&domain2.ListItem{}).Count(cursor.TotalCount.P).Error
	if e != nil {
		return out, e
	}
	db = db.Order("id " + strings.ToLower(string(order)))
	db = cursor.ApplyToGORM(db).Find(&out)
	if len(out) == 0 {
		return out, gorm.ErrRecordNotFound
	}
	return out, db.Error
}

func NewListRepo(db *gorm.DB) domain2.ListRepo {
	return &ListRepo{db: db}
}
