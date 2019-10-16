package eventbus

import (
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	stan "github.com/nats-io/stan.go"
	"github.com/qilin/crm-api/internal/eventbus/common"
)

type Publisher struct {
	conn       stan.Conn
	log        logger.Logger
	marshaller common.Marshaller
	subject    string
	wrapper    common.Wrapper
}

func (p *Publisher) Publish(msg common.Payloader, attempt int) error {
	evt, err := p.wrapper.Wrap(msg, attempt)
	if err != nil {
		p.log.Error("wrapping payload failed: %v", logger.Args(err))
		return err
	}
	return p.PublishEvent(evt)
}

func (b *Publisher) PublishEvent(evt common.Event) error {
	data, err := b.marshaller.Marshall(evt)
	if err != nil {
		b.log.Error("marshalling event failed: %v", logger.Args(err))
		return err
	}
	return b.conn.Publish(b.subject, data)
}

func NewPublisher(conn stan.Conn, log logger.Logger, subject string, marshaller common.Marshaller, wrapper common.Wrapper) *Publisher {
	return &Publisher{
		conn:       conn,
		log:        log,
		marshaller: marshaller,
		subject:    subject,
		wrapper:    wrapper,
	}
}
