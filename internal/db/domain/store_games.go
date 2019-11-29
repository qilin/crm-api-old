package domain

import (
	"context"
	"time"
)

// todo: probably rename
type StoreGamesItem struct {
	ID        string    `gorm:"column:id" json:"id"`
	URL       string    `gorm:"column:url" json:"url"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName
func (StoreGamesItem) TableName() string {
	// todo: rename, conflict with hasura table name
	return "products"
}

type StoreGamesRepo interface {
	All(ctx context.Context, limit, offset int) ([]StoreGamesItem, error)
	Create(ctx context.Context, model *StoreGamesItem) error
	Get(ctx context.Context, id string) (*StoreGamesItem, error)
	Delete(ctx context.Context, item *StoreGamesItem) error
}
