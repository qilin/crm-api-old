package eventbus

import (
	"context"
	"sync"

	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	stan "github.com/nats-io/stan.go"
	stan2 "github.com/qilin/crm-api/internal/stan"
)

type EventBus struct {
	mu   sync.Mutex
	ctx  context.Context
	cfg  *Config
	stan *stan2.Stan
	conn stan.Conn
	provider.LMT
	invites *Publisher
}

func (a *EventBus) Run() error {
	var e error
	a.conn, e = a.stan.SimpleConnect(a.ctx, stan2.ConnCfg{
		Backoff:      a.stan.BackoffCopy(),
		StanClientID: "qilin",
	}, a.L())
	if e != nil {
		a.L().WithFields(logger.Fields{"cmp": "eventbus"}).Error("simple connect error: %v", logger.Args(e.Error()))
		return e
	}

	a.L().Info("eventbus started")

	m := NewJSONMarshaller()
	a.invites = NewPublisher(a.conn, a.L().WithFields(logger.Fields{
		"publisher": "invites",
	}), "invites", NewJSONMarshaller(), NewJsonWrapper(m))

	a.conn.Subscribe("invites", func(msg *stan.Msg) {
		go func() {
			var (
				evt Event
				inv Invite
			)
			err := m.UnMarshall(msg.Data, &evt)
			if err != nil {
				a.L().Error("subscribe msg.Data unmarshal error: %v", logger.Args(err))
				return
			}
			err = m.UnMarshall(evt.Payload, &inv)
			if err != nil {
				a.L().Error("subscribe event.Payload unmarshal error: %v", logger.Args(err))
				return
			}
			a.L().Info("subscription invite: %v", logger.Args(inv))
		}()
	})

	return nil
}

func (a *EventBus) Invites() *Publisher {
	return a.invites
}

func (a *EventBus) Stop() error {
	return a.conn.Close()
}

// Config
type Config struct {
	Debug   bool `fallback:"shared.debug"`
	invoker *invoker.Invoker
}

// OnReload
func (c *Config) OnReload(callback func(ctx context.Context)) {
	c.invoker.OnReload(callback)
}

// Reload
func (c *Config) Reload(ctx context.Context) {
	c.invoker.Reload(ctx)
}

func New(ctx context.Context, set provider.AwareSet, stan *stan2.Stan, cfg *Config) *EventBus {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": Prefix})
	return &EventBus{
		ctx:  ctx,
		cfg:  cfg,
		stan: stan,
		LMT:  &set,
	}
}
