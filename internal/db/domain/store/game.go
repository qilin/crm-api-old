package store

type Game struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	Description string     `json:"description"`
	Publisher   *Publisher `json:"publisher"`
	Covers      *Covers    `json:"covers"`
	Screenshots []*Image   `json:"screenshots"`
	Tags        []*Tag     `json:"tags"`
	Genre       Genre      `json:"genre"`
	Rating      int        `json:"rating"`
}

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

type Image struct {
	URL string `json:"url"`
}

type Publisher struct {
	Title string `json:"title"`
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
