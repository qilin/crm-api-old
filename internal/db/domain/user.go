package domain

import (
	"context"
	"time"
)

type UserItem struct {
	ID            int       `gorm:"column:id" json:"id"`
	TenantID      int       `gorm:"column:tenant_id" json:"tenant_id"`
	Status        bool      `gorm:"column:status" json:"status"`
	Email         string    `gorm:"column:email" json:"email"`
	Picture       string    `gorm:"column:picture" json:"picture"`
	FirstName     string    `gorm:"column:first_name" json:"first_name"`
	LastName      string    `gorm:"column:last_name" json:"last_name"`
	Role          string    `gorm:"column:role" json:"role"`
	ExternalID    string    `gorm:"column:external_id" json:"external_id"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
	AuthTimestamp time.Time `gorm:"column:auth_timestamp" json:"-"`
	Password      string    `gorm:"-" json:"-"`
}

// TableName
func (UserItem) TableName() string {
	return "users"
}

type UserRepo interface {
	Create(ctx context.Context, model *UserItem) error
	Get(ctx context.Context, id int) (*UserItem, error)
	FindByEmail(ctx context.Context, email string) (*UserItem, error)
	FindByExternalID(ctx context.Context, externalID string) (*UserItem, error)
	IsExistsEmail(ctx context.Context, email string) (bool, error)
}
