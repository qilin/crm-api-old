package repo

import (
	"context"

	"github.com/qilin/crm-api/internal/db/trx"

	"github.com/jinzhu/gorm"
	"github.com/qilin/crm-api/internal/db/domain"
)

type AuthLog struct {
	db *gorm.DB
}

func (a *AuthLog) Create(ctx context.Context, model *domain.AuthLogItem) error {
	db := trx.Inject(ctx, a.db)
	return db.Save(model).Error
}

func (a *AuthLog) All(ctx context.Context, user_id int, offset, limit int) ([]domain.AuthLogItem, error) {
	db := trx.Inject(ctx, a.db)
	var (
		out = []domain.AuthLogItem{}
		e   error
	)
	e = db.Where("user_id=?", user_id).Offset(offset).Limit(limit).Order("created_at desc").Find(&out).Error
	return out, e
}

func NewAuthLogRepo(db *gorm.DB) domain.AuthLogRepo {
	return &AuthLog{db: db}
}

func AuthLogRepoProvider(db *gorm.DB) domain.AuthLogRepo {
	return NewAuthLogRepo(db)
}
