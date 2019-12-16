package store

import (
	"encoding/json"
	"errors"

	"github.com/tidwall/gjson"
)

type Game interface {
	Common() *GameCommon
}

type GameType string

const (
	GameTypeWeb     GameType = "Web"
	GameTypeDesktop GameType = "Desktop"
)

type GameCommon struct {
	ID   string   `json:"id"`
	Type GameType `json:"type"`
	Slug string   `json:"slug"`

	Title        string                `json:"title"`
	Summary      string                `json:"summary"`
	Description  string                `json:"description"`
	Developer    string                `json:"developer"`
	Publisher    string                `json:"publisher"`
	Genres       []Genre               `json:"genres"`
	ReleaseDate  *string               `json:"releaseDate"`
	Media        *Media                `json:"media"`
	Tags         []*Tag                `json:"tags"`
	Requirements []*SystemRequirements `json:"requirements"`
	Languages    *Languages            `json:"languages"`
	Platforms    []Platform            `json:"platforms"`
	Rating       int                   `json:"rating"`

	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Discount int     `json:"discount"`
}

func (g *GameCommon) Common() *GameCommon { return g }

type WebGame struct {
	GameCommon
}

type DesktopGame struct {
	GameCommon
}

type Platform string

const (
	PlatformLinux   Platform = "Linux"
	PlatformMacOs   Platform = "MacOS"
	PlatformWindows Platform = "Windows"
	PlatformWeb     Platform = "Web"
)

type Covers struct {
	FloorSmall    *Image `json:"floor_small"`
	MainLittle    *Image `json:"main_little"`
	MainBig       *Image `json:"main_big"`
	FloorMedium   *Image `json:"floor_medium"`
	FloorWide     *Image `json:"floor_wide"`
	FloorLarge    *Image `json:"floor_large"`
	FloorSmallest *Image `json:"floor_smallest"`
	FloorWidest   *Image `json:"floor_widest"`
	BackgroundBig *Image `json:"background_big"`
}

type SystemRequirements struct {
	Platform    string           `json:"platform"`
	Minimal     *RequirementsSet `json:"minimal"`
	Recommended *RequirementsSet `json:"recommended"`
}

type Media struct {
	Screenshots []*Image `json:"screenshots"`
	Trailers    []*Video `json:"trailers"`
}

type RequirementsSet struct {
	CPU       *string `json:"cpu"`
	DiskSpace *string `json:"diskSpace"`
	Gpu       *string `json:"gpu"`
	Os        *string `json:"os"`
	RAM       *string `json:"ram"`
}

type Languages struct {
	Audio []string `json:"audio"`
	Text  []string `json:"text"`
}

type Image struct {
	URL string `json:"url"`
}

type Video struct {
	URL string `json:"url"`
}

type Tag struct {
	Name string  `json:"name"`
	Type TagType `json:"type"`
}

type Genre string

const (
	GenreBoard     Genre = "Board"
	GenreCards     Genre = "Cards"
	GenreCasino    Genre = "Casino"
	GenreFarm      Genre = "Farm"
	GenreRacing    Genre = "Racing"
	GenreShooter   Genre = "Shooter"
	GenreFindItems Genre = "FindItems"
	GenrePuzzle    Genre = "Puzzle"
	GenreRpg       Genre = "RPG"
	GenreSimulator Genre = "Simulator"
	GenreStrategy  Genre = "Strategy"
)

type TagType string

const (
	TagTypeGenre  TagType = "genre"
	TagTypeCommon TagType = "common"
)

func UnmarshalGame(raw []byte) (Game, error) {
	t := GameType(gjson.GetBytes(raw, "type").String())
	return UnmarshalGameType(t, raw)
}

func UnmarshalGameType(t GameType, raw []byte) (Game, error) {
	switch t {
	case GameTypeWeb:
		var m WebGame
		return &m, json.Unmarshal(raw, &m)
	case GameTypeDesktop:
		var m DesktopGame
		return &m, json.Unmarshal(raw, &m)
	}
	return nil, errors.New("unknwon game type")
}

type GameSlice []Game

func (g *GameSlice) UnmarshalJSON(data []byte) error {
	var v []json.RawMessage
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	res := make([]Game, len(v))
	for i := range v {
		game, err := UnmarshalGame(v[i])
		if err != nil {
			return err
		}
		res[i] = game
	}

	*g = res
	return nil
}
