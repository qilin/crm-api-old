package sdk

type Request struct {
	URL  string `json:"url" form:"url" query:"url"`
	Meta string `json:"meta" form:"meta" query:"meta"`
}
