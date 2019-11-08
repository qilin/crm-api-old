package domain

import "context"

// todo: probably rename
type PlatformItem struct {
	ID     int    `gorm:"column:id" json:"id"`
	Name   string `gorm:"name", json:"name"`
	Status bool   `gorm:"column:status" json:"status"`
}

// TableName
func (PlatformItem) TableName() string {
	// todo: rename, conflict with hasura table name
	return "platforms"
}

type PlatformRepo interface {
	Create(ctx context.Context, model *PlatformItem) error
	Get(ctx context.Context, id int) (*PlatformItem, error)
}
