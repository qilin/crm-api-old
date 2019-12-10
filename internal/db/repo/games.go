package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/qilin/crm-api/internal/db/domain/store"
)

type history struct {
	ID      string
	Version int `gorm:"AUTO_INCREMENT"`
	Data    postgres.Jsonb
	Deleted bool
}

func (*history) TableName() string { return "store.games_history" }

type gameItem struct {
	ID      string
	Version int
	Data    postgres.Jsonb
}

func (*gameItem) TableName() string { return "store.games" }

type GamesRepo struct {
	db *gorm.DB
}

func NewGamesRepo(db *gorm.DB) *GamesRepo {
	return &GamesRepo{db}
}

func (r *GamesRepo) Insert(ctx context.Context, game *store.Game) (err error) {

	raw, err := json.Marshal(game)
	if err != nil {
		return err
	}
	// Note the use of tx as the database handle once you are within a transaction

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to store game: %s", r)
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	h := &history{ID: game.ID, Data: postgres.Jsonb{raw}}
	if err := tx.Create(h).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Save(&gameItem{ID: game.ID, Version: h.Version, Data: postgres.Jsonb{raw}}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *GamesRepo) Delete(ctx context.Context, id string) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	h := &history{ID: id, Deleted: true}
	if err := tx.Create(h).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&gameItem{ID: id}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *GamesRepo) Get(ctx context.Context, id string) (*store.Game, error) {
	var item gameItem
	if err := r.db.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}

	var g store.Game
	if err := json.Unmarshal(item.Data.RawMessage, &g); err != nil {
		return nil, err
	}

	return &g, nil
}

func (r *GamesRepo) All(ctx context.Context) ([]*store.Game, error) {
	var items []gameItem
	if err := r.db.Find(&items).Error; err != nil {
		return nil, err
	}

	var games = make([]*store.Game, 0, len(items))
	for i := range items {
		var g store.Game
		if err := json.Unmarshal(items[i].Data.RawMessage, &g); err != nil {
			return nil, err
		}
		games = append(games, &g)
	}
	return games, nil
}
