package eventbus

import (
	"context"

	"github.com/qilin/crm-api/internal/eventbus/common"

	"github.com/qilin/crm-api/internal/eventbus/subscribers"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/eventbus/publishers"
	"github.com/qilin/crm-api/internal/stan"
)

// Cfg
func Cfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{
		invoker: invoker.NewInvoker(),
	}
	e := cfg.UnmarshalKeyOnReload(UnmarshalKey, c)
	return c, func() {}, e
}

// CfgTest
func CfgTest() (*Config, func(), error) {
	c := &Config{
		invoker: invoker.NewInvoker(),
	}
	return c, func() {}, nil
}

// Provider
func Provider(ctx context.Context, set provider.AwareSet, pubs common.Publishers, subs common.Subscribers, stan *stan.Stan, cfg *Config, stanCfg *stan.Config) (*EventBus, func(), error) {
	g := New(ctx, set, stan, pubs, subs, cfg, stanCfg)
	err := g.Run()
	cleanup := func() {
		g.Stop()
	}
	return g, cleanup, err
}

var (
	WireSet = wire.NewSet(
		Provider,
		publishers.ProviderPublishers,
		subscribers.ProviderSubscribers,
		Cfg,
	)
	WireTestSet = wire.NewSet(
		Provider,
		publishers.ProviderPublishers,
		subscribers.ProviderSubscribers,
		CfgTest,
	)
)
