package sdk

type AuthRequest struct {
	URL  string `json:"url" form:"url" query:"url"`
	Meta string `json:"meta" form:"meta" query:"meta"`
}

type OrderRequest struct {
	Data string `json:"data" form:"data" query:"url"`
}
