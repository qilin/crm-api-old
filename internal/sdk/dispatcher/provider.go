package dispatcher

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/dispatcher/common"
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

func ProviderDispatcher(ctx context.Context, set provider.AwareSet, app AppSet, cfg *Config) (*Dispatcher, func(), error) {
	d := New(ctx, set, app, cfg)
	return d, func() {}, nil
}

var (
	WireSet = wire.NewSet(
		Cfg,
		ProviderDispatcher,
		wire.Struct(new(AppSet), "*"),
	)

	WireTestSet = wire.NewSet(
		CfgTest,
		ProviderDispatcher,
		wire.Struct(new(AppSet), "*"),
	)
)
