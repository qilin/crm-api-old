package domain

import (
	"context"
	"time"
)

type PlatformJWTKeyItem struct {
	ID         int       `gorm:"column:id" json:"id"`
	PlatformID int       `gorm:"platform_id" json:"platform_id"`
	Alg        string    `gorm:"column:alg" json:"alg"`
	Iss        string    `gorm:"column:iss" json:"alg"`
	KeyType    string    `gorm:"key_type" json:"key_type"`
	Key        string    `gorm:"column:key" json:"-"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName
func (PlatformJWTKeyItem) TableName() string {
	return "platform_jwt_keys"
}

type PlatformJWTKeyRepo interface {
	All(ctx context.Context) ([]PlatformJWTKeyItem, error)
	Create(ctx context.Context, model *PlatformJWTKeyItem) error
	Get(ctx context.Context, alg, iss string) (*PlatformJWTKeyItem, error)
	Delete(ctx context.Context, item *PlatformJWTKeyItem) error
}
