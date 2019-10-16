package eventbus

import (
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	stan "github.com/nats-io/stan.go"
)

type Publisher struct {
	conn       stan.Conn
	log        logger.Logger
	marshaller Marshaller
	subject    string
	wrapper    Wrapper
}

func (p *Publisher) Publish(msg Payloader, attempt int) error {
	evt, err := p.wrapper.Wrap(msg, attempt)
	if err != nil {
		p.log.Error("wrapping payload failed: %v", logger.Args(err))
		return err
	}
	return p.PublishEvent(evt)
}

func (b *Publisher) PublishEvent(evt Event) error {
	data, err := b.marshaller.Marshall(evt)
	if err != nil {
		b.log.Error("marshalling event failed: %v", logger.Args(err))
		return err
	}
	return b.conn.Publish("invites", data)
}

func NewPublisher(conn stan.Conn, log logger.Logger, subject string, marshaller Marshaller, wrapper Wrapper) *Publisher {
	return &Publisher{
		conn:       conn,
		log:        log,
		marshaller: marshaller,
		subject:    subject,
		wrapper:    wrapper,
	}
}
