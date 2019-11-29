package domain

import (
	"context"
	"time"
)

type UserMapItem struct {
	UserID     string    `gorm:"column:user_id" json:"user_id"`
	StoreID    int       `gorm:"column:store_id" json:"store_id"`
	ExternalID string    `gorm:"column:external_id" json:"external_id"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName
func (UserMapItem) TableName() string {
	return "user_map"
}

type UserMapRepo interface {
	Create(ctx context.Context, model *UserMapItem) error
	Get(ctx context.Context, id int) (*UserMapItem, error)
	FindByExternalID(ctx context.Context, storeId int, externalID string) (*UserMapItem, error)
	IsExists(ctx context.Context, id string) (bool, error)
}
