package common

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	stan "github.com/nats-io/stan.go"
)

const (
	Prefix       = "internal.eventbus"
	UnmarshalKey = "eventbus"
)

// Config: names for subjects
type Subjects struct {
	InvitesOut string
	InvitesIn  string
}

//
type SubscribeHandler func(ctx context.Context, sc stan.Conn, log logger.Logger) error

// EventBus
type EventBus interface {
	Publish(msg Payloader) error
	PublishEvent(evt Event) error
	Subscribe(SubscribeHandler) error
}

// Subscribers & Publishers
type Subscriber interface {
	Name() string
	Subscribe(conn stan.Conn, eb EventBus, subjects Subjects, l logger.Logger) error
	Unsubscribe() error
}

type Subscribers []Subscriber

type Publisher interface {
	Init(subs Subjects, conn stan.Conn, log logger.Logger)
	Name() string
	Publish(msg Payloader) error
	PublishEvent(evt Event) error
}

type Publishers []Publisher

// Event
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
