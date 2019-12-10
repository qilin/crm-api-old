package store

import (
	"encoding/json"
	"errors"
)

type Module interface {
	GetID() string
	GetType() ModuleType
	GetCategory() UserCategory
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
	Title        string       `json:"title"`
	Description  *string      `json:"description"`
	Link         string       `json:"link"`
	Image        *Image       `json:"image"`
	Backgound    *string      `json:"backgound"`
}

func (m *Breaker) GetID() string             { return m.ID }
func (m *Breaker) GetType() ModuleType       { return m.Type }
func (m *Breaker) GetCategory() UserCategory { return m.UserCategory }

type FreeGamesGroup struct {
	ID           string           `json:"id"`
	Type         ModuleType       `json:"type"`
	UserCategory UserCategory     `json:"user_category"`
	Title        string           `json:"title"`
	Games        []*FreeGameOffer `json:"games"`
}

func (m *FreeGamesGroup) GetID() string             { return m.ID }
func (m *FreeGamesGroup) GetType() ModuleType       { return m.Type }
func (m *FreeGamesGroup) GetCategory() UserCategory { return m.UserCategory }

type FreeGameOffer struct {
	Game  *Game  `json:"game"`
	Image *Image `json:"image"`
}

type StoreFront struct {
	Modules []Module `json:"modules"`
}

func UnmarshalModule(t ModuleType, raw []byte) (Module, error) {
	switch t {
	case ModuleTypeBreaker:
		var m Breaker
		return &m, json.Unmarshal(raw, &m)
	case ModuleTypeFreeGamesGroup:
		var m FreeGamesGroup
		return &m, json.Unmarshal(raw, &m)
	}
	return nil, errors.New("unknwon module type")
}
