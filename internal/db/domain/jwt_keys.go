package domain

import (
	"context"
	"time"
)

type JWTKeysItem struct {
	ID        int       `gorm:"column:id" json:"id"`
	Alg       string    `gorm:"column:alg" json:"alg"`
	Iss       string    `gorm:"column:iss" json:"alg"`
	KeyType   string    `gorm:"key_type" json:"key_type"`
	Key       string    `gorm:"column:key" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName
func (JWTKeysItem) TableName() string {
	return "jwt_keys"
}

type JWTKeysRepo interface {
	All(ctx context.Context) ([]JWTKeysItem, error)
	Create(ctx context.Context, model *JWTKeysItem) error
	Get(ctx context.Context, alg, iss string) (*JWTKeysItem, error)
	Delete(ctx context.Context, item *JWTKeysItem) error
}
