package sdk

type ErrorResponse struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type AuthResponse struct {
	Token string      `json:"token,omitempty"`
	Meta  interface{} `json:"meta,omitempty"`
}

type OrderResponse struct {
	// response
	Version int `json:"version"`
	// Merchant transaction if
	TransactionId string `json:"transaction_id"`
	//
	OrderId string `json:"order_id"`
	// Title
	Title string `json:"title"`
	// Text description
	Description string `json:"description"`
	// URL to image
	PreviewImage string `json:"preview_image"`
	// payment amount
	Amount string `json:"amount"`
	//
	Currency string `json:"currency"`
	// custom data
	Data Merchant `json:"data"`
	// timestamp
	Ts int `json:"ts"`
}

type Merchant struct {
	MerchantData string `json:"merchant_data"`
	// merchant_sing = sha1(merchant_data + merchant_secret_key)
	MerchantSign string `json:"merchant_sign"`
}
