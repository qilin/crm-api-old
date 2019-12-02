package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/qilin/crm-api/internal/db/domain"
	"github.com/qilin/crm-api/internal/db/trx"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type ActionsLog struct {
	db *gorm.DB
}

type record struct {
	Type      string      `json:"type"`
	Version   string      `json:"version"`
	Timestamp time.Time   `json:"timestamp`
	Data      interface{} `json:"data"`
}

func (a *ActionsLog) Add(ctx context.Context, act domain.Action) error {

	r := &record{
		Type:      act.Type(),
		Version:   act.Version(),
		Timestamp: time.Now().UTC(),
		Data:      act,
	}

	data, err := json.Marshal(r)
	if err != nil {
		return err
	}

	type GormWrapper struct {
		Data postgres.Jsonb
	}

	action := &GormWrapper{
		Data: postgres.Jsonb{json.RawMessage(data)},
	}
	db := trx.Inject(ctx, a.db)
	return db.Create(action).Error
}

func NewActionsLog(db *gorm.DB) *ActionsLog {
	return &ActionsLog{db: db}
}

func ActionsLogProvider(db *gorm.DB) domain.ActionsLog {
	return NewActionsLog(db)
}
