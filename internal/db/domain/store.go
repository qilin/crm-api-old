package domain

import "context"

// todo: probably rename
type StoreItem struct {
	ID       int    `gorm:"column:id" json:"id"`
	TenantId int    `gorm:"column:tenant_id" json:"tenant_id"`
	Name     string `gorm:"name", json:"name"`
	Status   bool   `gorm:"column:status" json:"status"`
}

// TableName
func (StoreItem) TableName() string {
	return "store"
}

type StoreRepo interface {
	Create(ctx context.Context, model *StoreItem) error
	Get(ctx context.Context, id int) (*StoreItem, error)
}
