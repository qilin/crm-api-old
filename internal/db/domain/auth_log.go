// module: Store Auth
package domain

import (
	"context"
	"time"
)

const (
	ActionSignIn  = "sign_in"
	ActionSignOut = "sign_out"
)

type AuthLogItem struct {
	ID int `gorm:"column:user_id" json:"id"`
	// user_id
	UserID int `gorm:"column:user_id" json:"id"`
	// log
	Action    string `gorm:"column:action" json:"action"`
	UserAgent string `gorm:"column:user_agent" json:"user_agent"`
	IP        string `gorm:"column:ip" json:"ip"`
	HWID      string `gorm:"column:hw_id" json:"hw_id"`
	// ts
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName
func (AuthLogItem) TableName() string {
	return "users.auth_log"
}

type AuthLogRepo interface {
	Create(ctx context.Context, model *AuthLogItem) error
	All(ctx context.Context, user_id int, offset, limit int) ([]AuthLogItem, error)
}
