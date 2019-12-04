// module: Store Auth
package domain

import (
	"context"
	"time"
)

type UsersItem struct {
	ID int `gorm:"column:id" json:"id"`

	Email    string `gorm:"column:email" json:"email"`
	Phone    string `gorm:"column:phone" json:"phone"`
	Password string `gorm:"column:password" json:"-"`

	Status       byte `gorm:"column:status" json:"status"`
	ServiceLevel byte `gorm:"column:service_level" json:"service_level"`

	Address1 string `gorm:"column:address_1" json:"address_1"`
	Address2 string `gorm:"column:address_2" json:"address_2"`
	City     string `gorm:"column:city" json:"city"`
	State    string `gorm:"column:state" json:"state"`
	Country  string `gorm:"column:country" json:"country"`
	Zip      string `gorm:"column:zip" json:"zip"`

	FirstName string `gorm:"column:first_name" json:"first_name"`
	LastName  string `gorm:"column:last_name" json:"last_name"`
	BirthDate int    `gorm:"column:birth_date" json:"birth_date"`

	Language string `gorm:"column:language" json:"language"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName
func (UsersItem) TableName() string {
	return "users.users"
}

type UsersRepo interface {
	Create(ctx context.Context, model *UsersItem) error
	Get(ctx context.Context, id int) (*UsersItem, error)
	FindByEmail(ctx context.Context, email string) (*UsersItem, error)
	FindByPhone(ctx context.Context, phone string) (*UsersItem, error)
	IsExistsEmail(ctx context.Context, email string) (bool, error)
	IsExistsPhone(ctx context.Context, phone string) (bool, error)
}
