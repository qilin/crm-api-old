package dispatcher

import (
	"context"

	"github.com/qilin/crm-api/internal/handlers"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/dispatcher/common"
)

// ProviderCfg
func ProviderCfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{
		invoker: invoker.NewInvoker(),
	}
	e := cfg.UnmarshalKeyOnReload(common.UnmarshalKey, c)
	return c, func() {}, e
}

// ProviderDispatcher
func ProviderDispatcher(ctx context.Context, set provider.AwareSet, appSet AppSet, cfg *Config) (*Dispatcher, func(), error) {
	d := New(ctx, set, appSet, cfg)
	return d, func() {}, nil
}

func ProviderSDKDispatcher(ctx context.Context, set provider.AwareSet, h common.Handlers, cfg *Config) (*SDKDispatcher, func(), error) {
	d := NewSDK(ctx, set, h, cfg)
	return d, func() {}, nil
}

var (
	WireSet = wire.NewSet(
		ProviderDispatcher,
		ProviderCfg,
		wire.Struct(new(AppSet), "*"),
	)

	WireTestSet = wire.NewSet(
		ProviderDispatcher,
		ProviderCfg,
		wire.Struct(new(AppSet), "*"),
	)

	SDKWireSet = wire.NewSet(
		ProviderSDKDispatcher,
		ProviderCfg,
		handlers.ProviderSDKHandlers,
	)

	SDKWireTestSet = wire.NewSet(
		ProviderSDKDispatcher,
		ProviderCfg,
		handlers.ProviderSDKHandlers,
	)
)
