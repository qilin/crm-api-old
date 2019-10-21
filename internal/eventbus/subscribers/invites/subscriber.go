package invites

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	stan "github.com/nats-io/stan.go"
	"github.com/qilin/crm-api/internal/eventbus/common"
	"github.com/qilin/crm-api/internal/eventbus/events"
)

type InviteSubscriber struct {
	cfg          *Config
	mailer       Mailer
	marshaller   common.Marshaller
	wrapper      common.Wrapper
	subscription stan.Subscription
}

func (s *InviteSubscriber) Subscribe(conn stan.Conn, eb common.EventBus, subs common.Subjects, log logger.Logger) error {
	var err error
	s.subscription, err = conn.Subscribe(subs.InvitesIn, func(msg *stan.Msg) {
		evt, err := s.wrapper.UnWrap(msg.Data)
		if err != nil {
			log.Error("can not unwrap event, error: %s", logger.Args(err.Error()))
			return
		}
		var invite events.Invite
		err = s.marshaller.UnMarshall(evt.Payload, &invite)
		if err != nil {
			log.Error("can not unmarshal event payload, error: %s", logger.Args(err.Error()))
			return
		}

		// todo: build message from template
		m := "you have new invite"
		if s.mailer.Send(invite.Email, s.cfg.Subject, m) != nil && evt.Attempt < s.cfg.MaxAttempts {
			evt.Attempt = evt.Attempt + 1
			eb.PublishEvent(evt)
		}
	})
	return err
}

func (s *InviteSubscriber) Unsubscribe() error {
	return s.subscription.Unsubscribe()
}

func (s *InviteSubscriber) Name() string {
	return "invites"
}

// Config
type Config struct {
	Debug   bool `fallback:"shared.debug"`
	Subject string
	Mailer  MailConfig
	// number of max attempts to send message on error
	MaxAttempts int
	invoker     *invoker.Invoker
}

type MailConfig struct {
	Host               string
	Port               int
	Username           string
	Password           string
	ReplyTo            string
	From               string
	InsecureSkipVerify bool
}

// OnReload
func (c *Config) OnReload(callback func(ctx context.Context)) {
	c.invoker.OnReload(callback)
}

// Reload
func (c *Config) Reload(ctx context.Context) {
	c.invoker.Reload(ctx)
}

func New(cfg *Config) *InviteSubscriber {
	m := common.NewJSONMarshaller()
	return &InviteSubscriber{
		cfg:        cfg,
		mailer:     NewMailer(cfg.Mailer),
		marshaller: m,
		wrapper:    common.NewJsonWrapper(m),
	}
}
