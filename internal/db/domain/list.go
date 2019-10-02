package domain

import (
	"context"

	"github.com/gurukami/typ/v2"
)

type ListItem struct {
	ID        typ.NullInt    `gorm:"column:id" json:"id"`
	Name      typ.NullString `gorm:"column:name" json:"name"`
	CreatedAt typ.NullTime   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt typ.NullTime   `gorm:"column:updated_at" json:"updated_at"`
}

// TableName
func (ListItem) TableName() string {
	return "list"
}

type ListRepo interface {
	Create(ctx context.Context, model *ListItem) error
	List(ctx context.Context, projection []string, cursor *Cursor, order Order, search string) ([]*ListItem, error)
}
