package domain

import (
	"context"
	"time"
)

// todo: probably rename
type ProductItem struct {
	ID        string    `gorm:"column:id" json:"id"`
	URL       string    `gorm:"column:url" json:"url"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName
func (ProductItem) TableName() string {
	// todo: rename, conflict with hasura table name
	return "products"
}

type ProductsRepo interface {
	All(ctx context.Context, limit, offset int) ([]ProductItem, error)
	Create(ctx context.Context, model *ProductItem) error
	Get(ctx context.Context, id string) (*ProductItem, error)
	Delete(ctx context.Context, item *ProductItem) error
}
