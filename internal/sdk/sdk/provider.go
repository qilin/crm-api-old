package sdk

import (
	"context"

	dispatcher "github.com/qilin/crm-api/internal/dispatcher/sdk"

	"github.com/qilin/crm-api/pkg/http"

	"github.com/ProtocolONE/go-core/v2/pkg/config"

	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/db/repo"
	"github.com/qilin/crm-api/internal/db/trx"
	"github.com/qilin/crm-api/internal/resolver"
	"github.com/qilin/crm-api/internal/sdk/common"
	sdkRepo "github.com/qilin/crm-api/internal/sdk/repo"
	"github.com/qilin/crm-api/pkg/postgres"
)

func Cfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{}
	e := cfg.UnmarshalKeyOnReload(common.UnmarshalKey, c)
	return c, func() {}, e
}

// CfgTest
func CfgTest() (*Config, func(), error) {
	return &Config{}, func() {}, nil
}

func PluginsCfg(cfg config.Configurator) (*PluginsConfig, func(), error) {
	c := &PluginsConfig{}
	e := cfg.UnmarshalKeyOnReload(common.UnmarshalKey, c)
	return c, func() {}, e
}

func PluginsCfgTest(cfg config.Configurator) (*PluginsConfig, func(), error) {
	c := &PluginsConfig{}
	return c, func() {}, nil
}

// Provider
func Provider(ctx context.Context, set provider.AwareSet, repo *sdkRepo.Repo, cfg *Config, pCfg *PluginsConfig, init config.Initial) (*SDK, func(), error) {
	g := New(ctx, set, repo, cfg, pCfg, init)
	return g, func() {}, nil
}

var (
	ProviderRepo = wire.NewSet(
		repo.NewPlatformRepo,
		repo.NewPlatformJWTKeyRepo,
		repo.NewProductsRepo,
		repo.NewUserMapRepo,
		repo.ActionsLogProvider,
		trx.NewTrxManager,
	)

	ProviderSDKRepo = wire.NewSet(
		ProviderRepo,
		wire.Struct(new(sdkRepo.Repo), "*"),
		postgres.WireSet,
	)

	ProviderSDKTestRepo = wire.NewSet(
		ProviderRepo,
		wire.Struct(new(sdkRepo.Repo), "*"),
		postgres.WireTestSet,
	)

	WireSet = wire.NewSet(
		Cfg,
		PluginsCfg,
		Provider,
		ProviderSDKRepo,
		resolver.ValidatorsProduction,
		dispatcher.WireSet,
		http.Provider,
		http.Cfg,
	)

	WireTestSet = wire.NewSet(
		CfgTest,
		PluginsCfgTest,
		Provider,
		ProviderSDKTestRepo,
		resolver.ValidatorsTest,
		dispatcher.WireTestSet,
		http.Provider,
		http.CfgTest,
	)
)
