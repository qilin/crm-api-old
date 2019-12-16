package cache

import (
	"context"
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/eko/gocache/cache"
	"github.com/google/wire"
)

// ProviderCfg returns configuration for Cache manager
func ProviderCfg(cfg config.Configurator) (*Config, func(), error) {
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

// Provider returns casbin.Enforcer instance with resolved dependencies
func Provider(ctx context.Context, set provider.AwareSet, cfg *Config) (*Cache, func(), error) {
	c, e := New(ctx, set, cfg)
	return c, func() {}, e
}

var (
	WireSet     = wire.NewSet(Provider, ProviderCfg, wire.Bind(new(cache.CacheInterface), new(*Cache)))
	WireTestSet = wire.NewSet(Provider, CfgTest, wire.Bind(new(cache.CacheInterface), new(*Cache)))
)
