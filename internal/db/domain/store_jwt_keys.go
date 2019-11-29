package domain

import (
	"context"
	"time"
)

type StoreJWTKeyItem struct {
	ID        int       `gorm:"column:id" json:"id"`
	Kid       string    `gorm:"column:kid" json:"kid"`
	StoreID   int       `gorm:"store_id" json:"store_id"`
	Alg       string    `gorm:"column:alg" json:"alg"`
	Iss       string    `gorm:"column:iss" json:"alg"`
	KeyType   string    `gorm:"key_type" json:"key_type"`
	Key       string    `gorm:"column:key" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName
func (StoreJWTKeyItem) TableName() string {
	return "store_jwt_keys"
}

type StoreJWTKeyRepo interface {
	All(ctx context.Context) ([]StoreJWTKeyItem, error)
	Create(ctx context.Context, model *StoreJWTKeyItem) error
	Get(ctx context.Context, alg, iss string) (*StoreJWTKeyItem, error)
	GetByIss(ctx context.Context, iss string) (*StoreJWTKeyItem, error)
	GetByKID(ctx context.Context, kid string) (*StoreJWTKeyItem, error)
	Delete(ctx context.Context, item *StoreJWTKeyItem) error
}
