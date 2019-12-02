package domain

import "context"

// logging actions

type Action interface {
	Type() string
	Version() string
}

type ActionsLog interface {
	Add(ctx context.Context, act Action) error
}

type BuyItemAction struct {
	UserID   string  `json:"user_id"`
	GameID   string  `json:"game_id"`
	ItemID   string  `json:"item_id"`
	StoreID  string  `json:"store_id"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

func (a *BuyItemAction) Type() string    { return "buy-item" }
func (a *BuyItemAction) Version() string { return "0" }
