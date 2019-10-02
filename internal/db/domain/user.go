package domain

import (
	"context"
	"time"
)

type UserItem struct {
	ID        int       `gorm:"column:id" json:"id"`
	Email     string    `gorm:"column:email" json:"email"`
	Password  string    `gorm:"column:password" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName
func (UserItem) TableName() string {
	return "users"
}

type UserRepo interface {
	Create(ctx context.Context, model *UserItem) error
	Get(ctx context.Context, email string, password string) (*UserItem, error)
	IsExistsEmail(ctx context.Context, email string) (bool, error)
}
