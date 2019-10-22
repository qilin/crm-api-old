package subscribers

import (
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	stan "github.com/nats-io/stan.go"
	"github.com/qilin/crm-api/internal/eventbus/common"
	"github.com/qilin/crm-api/internal/eventbus/events"
)

type inviteSubscriber struct {
	marshaller   common.Marshaller
	wrapper      common.Wrapper
	subscription stan.Subscription
}

func NewInviteSubscriber() common.Subscriber {
	m := common.NewJSONMarshaller()
	return &inviteSubscriber{
		marshaller: m,
		wrapper:    common.NewJsonWrapper(m),
	}
}

func (s *inviteSubscriber) Subscribe(conn stan.Conn, eb common.EventBus, subs common.Subjects, log logger.Logger) error {
	var err error
	s.subscription, err = conn.Subscribe(subs.InvitesOut, func(msg *stan.Msg) {
		evt, err := s.wrapper.UnWrap(msg.Data)
		if err != nil {
			log.Error("can not unwrap event, error: %s", logger.Args(err.Error()))
			return
		}
		var invite events.Invite
		err = s.marshaller.Unmarshal(evt.Payload, &invite)
		if err != nil {
			log.Error("can not unmarshal event payload, error: %s", logger.Args(err.Error()))
			return
		}

		log.Info("Invites subscriber received invite: %v", logger.Args(invite))
	})
	return err
}

func (s *inviteSubscriber) Unsubscribe() error {
	return s.subscription.Unsubscribe()
}

func (s *inviteSubscriber) Name() string {
	return "invites"
}
