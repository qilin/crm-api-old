package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/qilin/crm-api/internal/db/domain/store"
)

type moduleHistory struct {
	ID           string
	UserCategory store.UserCategory
	Type         store.ModuleType
	Version      int `gorm:"AUTO_INCREMENT"`
	Data         postgres.Jsonb
	Deleted      bool
}

func (*moduleHistory) TableName() string { return "store.modules_history" }

type module struct {
	ID           string
	UserCategory store.UserCategory
	Type         store.ModuleType
	Version      int
	Data         postgres.Jsonb
}

func (*module) TableName() string { return "store.modules" }

type StorefrontRepo struct {
	db *gorm.DB
}

func NewStorefrontRepo(db *gorm.DB) *StorefrontRepo {
	return &StorefrontRepo{db}
}

func (r *StorefrontRepo) InsertModule(ctx context.Context, mod store.Module) (err error) {

	raw, err := json.Marshal(mod)
	if err != nil {
		return err
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to store module: %s", r)
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	t, err := r.getType(tx, mod.GetID())
	switch {
	case gorm.IsRecordNotFoundError(err):
		// fine
	case err != nil:
		return err
	case t != mod.GetType():
		return fmt.Errorf("Invalid module type '%s', expected '%s'", mod.GetType(), t)
	}

	h := &moduleHistory{
		ID:           mod.GetID(),
		UserCategory: mod.GetCategory(),
		Type:         mod.GetType(),
		Data:         postgres.Jsonb{raw},
	}
	if err := tx.Create(h).Error; err != nil {
		tx.Rollback()
		return err
	}

	m := &module{
		ID:           mod.GetID(),
		UserCategory: mod.GetCategory(),
		Type:         mod.GetType(),
		Version:      h.Version,
		Data:         postgres.Jsonb{raw},
	}
	if err := tx.Save(m).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *StorefrontRepo) getType(tx *gorm.DB, id string) (store.ModuleType, error) {
	var m module
	res := tx.Select("type").Where("id = ?", id).First(&m)
	return m.Type, res.Error
}

func (r *StorefrontRepo) Delete(ctx context.Context, id string, category store.UserCategory) (err error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to delete module: %s", r)
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	h := &moduleHistory{
		ID:           id,
		UserCategory: category,
		Deleted:      true,
	}
	if err := tx.Create(h).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&module{ID: id, UserCategory: category}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *StorefrontRepo) GetModule(ctx context.Context, id string, category store.UserCategory) (store.Module, error) {
	var m module
	if err := r.db.Where("id = ? and user_category = ?", id, category).First(&m).Error; err != nil {
		return nil, err
	}

	return store.UnmarshalModule(m.Type, m.Data.RawMessage)
}
