package store

import (
	"encoding/json"
	"errors"
)

type Module interface {
	GetID() string
	GetType() ModuleType
	GetCategory() UserCategory
	GetVersion() int
}

type ModuleType string

const (
	ModuleTypeBreaker        ModuleType = "Breaker"
	ModuleTypeFreeGamesGroup ModuleType = "FreeGamesGroup"
)

type UserCategory string

// predefined user category
const (
	UserCategoryUnknown = ""
)

type Breaker struct {
	ID           string       `json:"id"`
	Type         ModuleType   `json:"type"`
	UserCategory UserCategory `json:"user_category"`
	Version      int          `json:"version"`
	Title        string       `json:"title"`
	Description  *string      `json:"description"`
	Link         string       `json:"link"`
	Image        *Image       `json:"image"`
	Backgound    *string      `json:"backgound"`
}

func (m *Breaker) GetID() string             { return m.ID }
func (m *Breaker) GetType() ModuleType       { return m.Type }
func (m *Breaker) GetCategory() UserCategory { return m.UserCategory }
func (m *Breaker) GetVersion() int           { return m.Version }

type FreeGamesGroup struct {
	ID           string           `json:"id"`
	Type         ModuleType       `json:"type"`
	UserCategory UserCategory     `json:"user_category"`
	Version      int              `json:"version"`
	Title        string           `json:"title"`
	Games        []*FreeGameOffer `json:"games"`
}

func (m *FreeGamesGroup) GetID() string             { return m.ID }
func (m *FreeGamesGroup) GetType() ModuleType       { return m.Type }
func (m *FreeGamesGroup) GetCategory() UserCategory { return m.UserCategory }
func (m *FreeGamesGroup) GetVersion() int           { return m.Version }

type FreeGameOffer struct {
	GameID string `json:"game_id"`
	Image  *Image `json:"image"`
}

type StoreFront struct {
	Modules []Module `json:"modules"`
}

func UnmarshalModule(t ModuleType, version int, raw []byte) (Module, error) {
	switch t {
	case ModuleTypeBreaker:
		var m Breaker
		err := json.Unmarshal(raw, &m)
		m.Version = version
		return &m, err
	case ModuleTypeFreeGamesGroup:
		var m FreeGamesGroup
		err := json.Unmarshal(raw, &m)
		m.Version = version
		return &m, err
	}
	return nil, errors.New("unknwon module type")
}
