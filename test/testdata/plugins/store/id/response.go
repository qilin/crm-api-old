package id

type Response struct {
	RPC    string `json:"jsonrpc"`
	Result struct {
		Info struct {
			Email string `json:"email"`
		} `json:"info"`
		Profile struct {
			Firstname string `json:"firstname"`
			Lastname  string `json:"lastname"`
			Gender    string `json:"gender"`
			Birthdate int    `json:"birthdate"`
		} `json:"profile"`
		Status string `json:"status"`
	} `json:"result"`
}
