package sdk

type ErrorResponse struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type AuthResponse struct {
	Token string `json:"token,omitempty"`
	Meta  string `json:"meta,omitempty"`
}

type OrderResponse struct {
	Data string `json:"data"`
}
