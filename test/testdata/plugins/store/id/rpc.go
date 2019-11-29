package id

type RPCRequest struct {
	RPC    string              `json:"jsonrpc"`
	Method string              `json:"method"`
	Params []map[string]string `json:"params"`
}

type GetProfileInfoResponse struct {
	RPC    string `json:"jsonrpc"`
	Result struct {
		Info struct {
			Email   string `json:"email"`
			ChainId struct {
				Default string `json:"default"`
			} `json:"chain_id"`
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
