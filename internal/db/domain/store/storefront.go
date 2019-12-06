package store

type Module interface {
	IsModule()
}

type ModuleType string

const (
	ModuleTypeBreaker        ModuleType = "Breaker"
	ModuleTypeFreeGamesGroup ModuleType = "FreeGamesGroup"
)

type Breaker struct {
	Type        ModuleType `json:"type"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Link        string     `json:"link"`
	Image       *Image     `json:"image"`
	Backgound   *string    `json:"backgound"`
}

func (Breaker) IsModule() {}

type FreeGamesGroup struct {
	Type  ModuleType       `json:"type"`
	Title string           `json:"title"`
	Games []*FreeGameOffer `json:"games"`
}

func (FreeGamesGroup) IsModule() {}

type FreeGameOffer struct {
	Game  *Game  `json:"game"`
	Image *Image `json:"image"`
}

type StoreFront struct {
	Modules []Module `json:"modules"`
}
