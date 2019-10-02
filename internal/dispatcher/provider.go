package dispatcher

import (
	"context"

	"github.com/qilin/crm-api/internal/dispatcher/common"

	"github.com/google/wire"
	"github.com/qilin/go-core/config"
	"github.com/qilin/go-core/invoker"
	"github.com/qilin/go-core/provider"
)

// ProviderCfg
func ProviderCfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{
		WorkDir: cfg.WorkDir(),
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
)
