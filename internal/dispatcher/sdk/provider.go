package sdk

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/dispatcher/common"
	"github.com/qilin/crm-api/internal/handlers"
)

// ProviderCfg
func Cfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{
		invoker: invoker.NewInvoker(),
	}
	e := cfg.UnmarshalKeyOnReload(common.UnmarshalKey, c)
	return c, func() {}, e
}

// CfgTest
func CfgTest() (*Config, func(), error) {
	return &Config{}, func() {}, nil
}

func ProviderDispatcher(ctx context.Context, set provider.AwareSet, h common.Handlers, cfg *Config) (*Dispatcher, func(), error) {
	d := New(ctx, set, h, cfg)
	return d, func() {}, nil
}

var (
	WireSet = wire.NewSet(
		Cfg,
		ProviderDispatcher,
		handlers.ProviderSDKHandlers,
	)

	WireTestSet = wire.NewSet(
		CfgTest,
		ProviderDispatcher,
		handlers.ProviderSDKHandlers,
	)
)
