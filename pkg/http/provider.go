package http

import (
	"context"

	"github.com/qilin/crm-api/pkg/graphql"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/dispatcher"
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
func Provider(ctx context.Context, set provider.AwareSet, dispatcher Dispatcher, cfg *Config) (*HTTP, func(), error) {
	http := New(ctx, set, dispatcher, cfg)
	return http, func() {}, nil
}

var (
	WireSet = wire.NewSet(
		Provider,
		Cfg,
		dispatcher.WireSet,
		graphql.WireSet,
	)
	WireTestSet = wire.NewSet(
		Provider,
		CfgTest,
		dispatcher.WireSet,
		graphql.WireTestSet,
	)
)
