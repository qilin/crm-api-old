package eventbus

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

type Invite struct {
}
