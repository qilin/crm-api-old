package authentication

import (
	"context"

	"github.com/qilin/crm-api/internal/db/repo"

	"github.com/qilin/crm-api/internal/plugins"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
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
func Provider(ctx context.Context, pm *plugins.PluginManager, set provider.AwareSet, appSet AppSet, cfg *Config) (*AuthenticationService, func(), error) {
	g, e := New(ctx, pm, set, appSet, cfg)
	return g, func() {}, e
}

var (
	RepoProvider = wire.NewSet(
		repo.NewAuthLogRepo,
		repo.NewAuthProviderRepo,
		repo.NewUsersRepo,
		repo.NewUserMapRepo,
	)

	WireSet = wire.NewSet(
		RepoProvider,
		Provider,
		Cfg,
		wire.Struct(new(AppSet), "*"),
		plugins.WireSet,
	)
	WireTestSet = wire.NewSet(
		RepoProvider,
		Provider,
		CfgTest,
		wire.Struct(new(AppSet), "*"),
		plugins.WireTestSet,
	)
)
