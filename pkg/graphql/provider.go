package graphql

import (
	"context"
	"net/http"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/generated/graphql"
	"github.com/qilin/crm-api/internal/resolver"
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
