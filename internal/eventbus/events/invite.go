package events

import "time"

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
