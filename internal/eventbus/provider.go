package eventbus

import (
	"context"

	"github.com/qilin/crm-api/internal/stan"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
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
func Provider(ctx context.Context, set provider.AwareSet, stan *stan.Stan, cfg *Config) (*EventBus, func(), error) {
	g := New(ctx, set, stan, cfg)
	err := g.Run()
	cleanup := func() {
		g.Stop()
	}
	return g, cleanup, err
}

var (
	WireSet     = wire.NewSet(Provider, Cfg)
	WireTestSet = wire.NewSet(Provider, CfgTest)
)
