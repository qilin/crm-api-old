package sdk

type Response struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Meta string `json:"meta,omitempty"`
}
