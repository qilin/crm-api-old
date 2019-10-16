package eventbus

import (
	"context"
	"sync"

	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	stan "github.com/nats-io/stan.go"
	"github.com/qilin/crm-api/internal/eventbus/common"
	stan2 "github.com/qilin/crm-api/internal/stan"
)

type EventBus struct {
	mu   sync.Mutex
	ctx  context.Context
	cfg  *Config
	stan *stan2.Stan
	conn stan.Conn
	provider.LMT
	pubs    map[string]*Publisher
	invites *Publisher
	subs    common.Subscribers
}

func (a *EventBus) Run() error {
	var e error
	a.conn, e = a.stan.SimpleConnect(a.ctx, stan2.ConnCfg{
		Backoff: a.stan.BackoffCopy(),
	}, a.L())
	if e != nil {
		a.L().WithFields(logger.Fields{"cmp": "eventbus"}).Error("simple connect error: %v", logger.Args(e.Error()))
		return e
	}

	for _, s := range a.subs {
		e := s.Subscribe(a.conn, a, a.cfg.Subjects, a.L().WithFields(logger.Fields{"subscriber": s.Name()}))
		if e != nil {
			a.L().Error("subscription %s failed with error %s", logger.Args(s.Name(), e.Error()))
		} else {
			a.L().Info("subscription %s started", logger.Args(s.Name()))
		}
	}

	m := common.NewJSONMarshaller()
	a.invites = NewPublisher(a.conn, a.L().WithFields(logger.Fields{
		"publisher": "invites",
	}), a.cfg.Subjects.InvitesOut, common.NewJSONMarshaller(), common.NewJsonWrapper(m))

	return nil
}

func (a *EventBus) Publish(msg common.Payloader) error {
	return nil
}

func (a *EventBus) PublishEvent(evt common.Event) error {
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
	Debug    bool `fallback:"shared.debug"`
	Subjects common.Subjects
	invoker  *invoker.Invoker
}

// OnReload
func (c *Config) OnReload(callback func(ctx context.Context)) {
	c.invoker.OnReload(callback)
}

// Reload
func (c *Config) Reload(ctx context.Context) {
	c.invoker.Reload(ctx)
}

func New(ctx context.Context, set provider.AwareSet, stan *stan2.Stan, pubs common.Publishers, subs common.Subscribers, cfg *Config) *EventBus {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": Prefix})
	return &EventBus{
		ctx:  ctx,
		cfg:  cfg,
		stan: stan,
		//publishers: pubs,
		subs: subs,
		LMT:  &set,
	}
}
