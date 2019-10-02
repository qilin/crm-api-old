package domain

import (
	"strings"

	"github.com/gurukami/typ/v2"
	"github.com/jinzhu/gorm"
	graphql1 "github.com/qilin/crm-api/generated/graphql"
)

// Cursor represent pagination info
type Cursor struct {
	TotalCount  typ.NullInt    `sql:"-"`
	Limit       typ.NullInt    `sql:"-"`
	Offset      typ.NullInt    `sql:"-"`
	HasNextPage typ.NullBool   `sql:"-"`
	Cursor      typ.NullString `sql:"-"`
}

// Init
func (c *Cursor) Init() {
	if !c.TotalCount.Present() {
		c.TotalCount.Set(0)
	}
}

// CalculateHasNextPage determines whether there is a next page
func (c *Cursor) CalculateHasNextPage() {
	c.HasNextPage.Set(c.TotalCount.V()-(c.Limit.V()+c.Offset.V()) > 0)
}

// FromGQLCursor
func (c *Cursor) ApplyGQLCursor(cursor graphql1.CursorIn, maxLimit int) {
	limit := cursor.Limit
	if limit > maxLimit {
		limit = maxLimit
	}
	c.Limit.Set(limit)
	c.Offset.Set(cursor.Offset)
}

// ToGQLCursor
func (c *Cursor) ToGQLCursor() *graphql1.CursorOut {
	c.CalculateHasNextPage()
	return &graphql1.CursorOut{
		Count:  c.TotalCount.V(),
		Limit:  c.Limit.V(),
		Offset: c.Offset.V(),
		IsEnd:  !c.HasNextPage.V(),
		Cursor: c.Cursor.V(),
	}
}

// ApplyToGORM
func (c *Cursor) ApplyToGORM(db *gorm.DB) *gorm.DB {
	if c.Limit.V() > 0 {
		db = db.Limit(c.Limit.V())
	}
	if c.Offset.V() > 0 {
		db = db.Offset(c.Offset.V())
	}
	return db
}

// AsSQL
func (c *Cursor) AsSQL() string {
	var sql []string
	if c.Limit.V() > 0 {
		sql = append(sql, "LIMIT "+c.Limit.Typ().String().V())
	}
	if c.Offset.V() > 0 {
		sql = append(sql, "OFFSET "+c.Offset.Typ().String().V())
	}
	return " " + strings.Join(sql, " ") + " "
}

type Order string

const (
	OrderAsc  Order = "ASC"
	OrderDesc Order = "DESC"
)

// ToOrder
func ToOrder(order string) Order {
	switch {
	case strings.EqualFold(order, "desc"):
		return OrderDesc
	}
	return OrderAsc
}
