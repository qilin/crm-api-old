package eventbus

import "time"

const (
	EventVersion = "1.0"
)

type Event struct {
	// number of attempts
	Attempt int `json:"attempt"`
	// name of payload type
	Name string `json:"name"`
	// event data
	Payload []byte `json:"payload"`
	// version of Event structure
	Version string `json:"version"`
}

type Payloader interface {
	Name() string
	Payload() interface{}
}

type Invite struct {
	Id         int       `json:"id"`
	TenantId   int       `json:"tenant_id"`
	UserId     int       `json:"user_id"`
	Expiration time.Time `json:"expiration"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Accepted   bool      `json:"accepted"`
}

func (i Invite) Name() string {
	return "Invite"
}

func (i Invite) Payload() interface{} {
	return i
}
