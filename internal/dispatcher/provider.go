package dispatcher

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/dispatcher/common"
	"github.com/qilin/crm-api/internal/validators"
	"gopkg.in/go-playground/validator.v9"
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

// Validators
func ProviderValidators(v *validators.ValidatorSet) (validate *validator.Validate, _ func(), err error) {
	validate = validator.New()

	// add needed validators

	return validate, func() {}, nil
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
