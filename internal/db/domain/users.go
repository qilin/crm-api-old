// module: Store Auth
package domain

import (
	"context"
	"time"
)

type UserStatus int

func (u UserStatus) Int8() int8 {
	return int8(u)
}

const (
	// todo: define needed statuses
	UserActive UserStatus = iota + 1
	UserBlocked
	UserLocked
)

type UsersItem struct {
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
func (UsersItem) TableName() string {
	return "users"
}

type UsersRepo interface {
	Create(ctx context.Context, model *UsersItem) error
	Get(ctx context.Context, id int) (*UsersItem, error)
	FindByEmail(ctx context.Context, email string) (*UsersItem, error)
	FindByPhone(ctx context.Context, phone string) (*UsersItem, error)
	IsExistsEmail(ctx context.Context, email string) (bool, error)
	IsExistsPhone(ctx context.Context, phone string) (bool, error)
}
