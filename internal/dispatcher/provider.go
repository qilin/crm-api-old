package dispatcher

import (
	"context"

	jwtverifier "github.com/ProtocolONE/authone-jwt-verifier-golang"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/dispatcher/common"
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

func ProviderAuthCfg(cfg config.Configurator) (*common.AuthConfig, func(), error) {
	c := &common.AuthConfig{}
	e := cfg.UnmarshalKey(common.UnmarshalAuthConfigKey, c)
	return c, func() {}, e
}

// ProviderDispatcher
func ProviderDispatcher(ctx context.Context, set provider.AwareSet, appSet AppSet, cfg *Config, authCfg *common.AuthConfig) (*Dispatcher, func(), error) {
	d := New(ctx, set, appSet, cfg, authCfg)
	return d, func() {}, nil
}

// jwt verifier
func ProviderJwtVerifier(cfg *common.AuthConfig) *jwtverifier.JwtVerifier {
	return jwtverifier.NewJwtVerifier(jwtverifier.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		Scopes:       []string{"openid", "offline"},
		RedirectURL:  cfg.RedirectUrl,
		Issuer:       cfg.Issuer,
	})
}

var (
	WireSet = wire.NewSet(
		ProviderDispatcher,
		//ProviderValidators,
		ProviderAuthCfg,
		ProviderJwtVerifier,
		ProviderCfg,
		wire.Struct(new(AppSet), "*"),
	)

	WireTestSet = wire.NewSet(
		ProviderDispatcher,
		//ProviderValidators,
		ProviderAuthCfg,
		ProviderJwtVerifier,
		ProviderCfg,
		wire.Struct(new(AppSet), "*"),
	)
)
