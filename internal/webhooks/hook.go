package webhooks

import "time"

type Hook struct {
	Id    string `json:"id"`
	Table struct {
		Schema string `json:"schema"`
		Name   string `json:"name"`
	} `json:"table"`
	Trigger struct {
		Name string `json:"name"`
	} `json:"trigger"`
	Event struct {
		SessionVariable map[string]string `json:"session_variable"`
		Op              string            `json:"op"`
		Data            struct {
			Old map[string]string      `json:"old"`
			New map[string]interface{} `json:"new"`
		} `json:"data"`
	} `json:"event"`
	DeliveryInfo struct {
		CurrentRetry int `json:"current_retry"`
		MaxRetries   int `json:"max_retries"`
	}
	CreatedAt time.Time `json:"created_at"`
}
