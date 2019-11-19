package common

type AuthRequest struct {
	URL              string      `json:"url" form:"url" query:"url" validate:"required,uri"`
	Meta             interface{} `json:"meta" form:"meta" query:"meta"`
	QilinProductUUID string      `json:"qilinProductUUID" from:"qilinProductUUID" query:"qilinProductUUID" validate:"omitempty,uuid4"`
}

type OrderRequest struct {
	UserID string `json:"user_id" form:"user_id" query:"user_id"`
	ItemID string `json:"item_id" form:"item_id" query:"item_id"`
	GameID string `json:"game_id" form:"game_id" query:"game_id"`
	// Currency string `json:"currency" form:"currency" query:"currency"`
	// Data     string `json:"data" form:"data" query:"url"`
}
