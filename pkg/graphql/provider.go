package graphql

import (
	"context"
	"net/http"

	"github.com/google/wire"
	"github.com/qilin/crm-api/generated/graphql"
	"github.com/qilin/crm-api/pkg/resolver"
	"github.com/qilin/go-core/config"
	"github.com/qilin/go-core/invoker"
	"github.com/qilin/go-core/provider"
)

// Cfg
func Cfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{
		invoker: invoker.NewInvoker(),
	}
	e := cfg.UnmarshalKeyOnReload(UnmarshalKey, c)
	c.Middleware = []func(http.Handler) http.Handler{}
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
func Provider(ctx context.Context, resolver graphql.Config, set provider.AwareSet, cfg *Config) (*GraphQL, func(), error) {
	g := New(ctx, resolver, set, cfg)
	return g, func() {}, nil
}

var (
	WireSet     = wire.NewSet(Provider, Cfg, resolver.WireSet)
	WireTestSet = wire.NewSet(Provider, CfgTest, resolver.WireTestSet)
)
