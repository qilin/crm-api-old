// module: Store Auth
package domain

import (
	"context"
	"time"
)

type UserProviderMapItem struct {
	// user_id
	UserID int `gorm:"column:user_id" json:"id"`
	// provider
	Provider    string `gorm:"column:provider" json:"provider"`
	ProviderID  string `gorm:"column:provider_id" json:"provider_id"`
	ProviderKey string `gorm:"column:provider_key" json:"provider_key"`
	// ts
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName
func (UserProviderMapItem) TableName() string {
	return "users.authentication_providers"
}

type UserProviderMapRepo interface {
	Create(ctx context.Context, model *UserProviderMapItem) error
	Delete(ctx context.Context, model *UserProviderMapItem) error
	Get(ctx context.Context, user_id int, provider, provider_id string) (*UserProviderMapItem, error)
	GetByUserId(ctx context.Context, user_id int) ([]UserProviderMapItem, error)
}
