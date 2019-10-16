package publishers

import (
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	stan "github.com/nats-io/stan.go"
	"github.com/qilin/crm-api/internal/eventbus/common"
	"github.com/qilin/crm-api/internal/eventbus/events"
)

type invitesPublisher struct {
	conn       stan.Conn
	log        logger.Logger
	marshaller common.Marshaller
	subject    string
	wrapper    common.Wrapper
}

func (p *invitesPublisher) Publish(msg common.Payloader) error {
	evt, err := p.wrapper.Wrap(msg, 0)
	if err != nil {
		p.log.Error("wrapping payload failed: %v", logger.Args(err))
		return err
	}
	return p.PublishEvent(evt)
}

func (b *invitesPublisher) PublishEvent(evt common.Event) error {
	data, err := b.marshaller.Marshall(evt)
	if err != nil {
		b.log.Error("marshalling event failed: %v", logger.Args(err))
		return err
	}
	return b.conn.Publish(b.subject, data)
}

func (b *invitesPublisher) Init(subjects common.Subjects, conn stan.Conn, log logger.Logger) {
	b.subject = subjects.InvitesOut
	b.conn = conn
	b.log = log
}

func (b *invitesPublisher) Name() string {
	return events.Invite{}.Name()
}

func NewInvitesPublisher() common.Publisher {
	m := common.NewJSONMarshaller()
	return &invitesPublisher{
		marshaller: m,
		wrapper:    common.NewJsonWrapper(m),
	}
}
