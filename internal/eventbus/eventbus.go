package eventbus

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	stan "github.com/nats-io/stan.go"
	"github.com/qilin/crm-api/internal/eventbus/common"
	stan2 "github.com/qilin/crm-api/internal/stan"
)

type EventBus struct {
	ctx  context.Context
	cfg  *Config
	stan *stan2.Stan
	conn stan.Conn
	provider.LMT
	pubs map[string]common.Publisher
	subs common.Subscribers
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

	for _, p := range a.pubs {
		p.Init(a.cfg.Subjects, a.conn, a.L().WithFields(logger.Fields{"publisher": p.Name()}))
		a.pubs[p.Name()] = p
	}

	for _, s := range a.subs {
		e := s.Subscribe(a.conn, a, a.cfg.Subjects, a.L().WithFields(logger.Fields{"subscriber": s.Name()}))
		if e != nil {
			a.L().Error("subscription %s failed with error %s", logger.Args(s.Name(), e.Error()))
		} else {
			a.L().Info("subscription %s started", logger.Args(s.Name()))
		}
	}

	return nil
}

func (a *EventBus) Publish(msg common.Payloader) error {
	p, ok := a.pubs[msg.Name()]
	if !ok {
		a.L().Error("Publisher not found " + msg.Name())
		return nil
	}
	return p.Publish(msg)
}

func (a *EventBus) PublishEvent(evt common.Event) error {
	p, ok := a.pubs[evt.Name]
	if !ok {
		a.L().Error("Publisher not found %s", logger.Args(evt.Name))
		return nil
	}
	return p.PublishEvent(evt)
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
	publishers := map[string]common.Publisher{}
	for _, p := range pubs {
		publishers[p.Name()] = p
	}
	return &EventBus{
		ctx:  ctx,
		cfg:  cfg,
		stan: stan,
		pubs: publishers,
		subs: subs,
		LMT:  &set,
	}
}
