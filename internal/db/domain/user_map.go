package domain

import (
	"context"
	"time"
)

type UserMapItem struct {
	UserID     string    `gorm:"column:id" json:"id"`
	ExternalID string    `gorm:"column:external_id" json:"external_id"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName
func (UserMapItem) TableName() string {
	return "users"
}

type UserMapRepo interface {
	Create(ctx context.Context, model *UserMapItem) error
	Get(ctx context.Context, id int) (*UserMapItem, error)
	FindByExternalID(ctx context.Context, externalID string) (*UserMapItem, error)
	IsExists(ctx context.Context, id string) (bool, error)
}
