package authentication

import (
	"context"

	"github.com/qilin/crm-api/internal/authentication/providers"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/db/repo"
	"github.com/qilin/crm-api/internal/plugins"
	"gopkg.in/go-playground/validator.v9"
)

// Cfg
func Cfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{}
	e := cfg.UnmarshalKey(UnmarshalKey, c)
	return c, func() {}, e
}

// CfgTest
func CfgTest() (*Config, func(), error) {
	c := &Config{}
	return c, func() {}, nil
}

// Provider
func Provider(ctx context.Context, pm *plugins.PluginManager, set provider.AwareSet, appSet AppSet, validate *validator.Validate, p1 *providers.P1Provider, cfg *Config) (*AuthenticationService, func(), error) {
	g, e := New(ctx, pm, set, appSet, validate, p1, cfg)
	return g, func() {}, e
}

var (
	RepoProvider = wire.NewSet(
		repo.NewAuthLogRepo,
		repo.NewAuthProviderRepo,
		//repo.NewUsersRepo,
		//repo.NewUserMapRepo,
	)

	WireSet = wire.NewSet(
		providers.WireSet,
		RepoProvider,
		Provider,
		Cfg,
		wire.Struct(new(AppSet), "*"),
	)
	WireTestSet = wire.NewSet(
		providers.WireSet,
		RepoProvider,
		Provider,
		CfgTest,
		wire.Struct(new(AppSet), "*"),
	)
)
