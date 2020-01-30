package sdk

import (
	"context"

	"github.com/qilin/crm-api/internal/authentication"
	dispatcher2 "github.com/qilin/crm-api/internal/sdk/dispatcher"

	"github.com/qilin/crm-api/internal/plugins"

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

// Provider
func Provider(ctx context.Context, pm *plugins.PluginManager, set provider.AwareSet, repo *sdkRepo.Repo, cfg *Config) (*SDK, func(), error) {
	g := New(ctx, pm, set, repo, cfg)
	return g, func() {}, nil
}

var (
	ProviderRepo = wire.NewSet(
		repo.NewPlatformRepo,
		repo.NewPlatformJWTKeyRepo,
		repo.NewProductsRepo,
		repo.NewUserMapRepo,
		repo.NewUsersRepo,
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
		Provider,
		ProviderSDKRepo,
		resolver.ValidatorsProduction,
		dispatcher2.WireSet,
		plugins.WireSet,
		authentication.WireSet,
		http.Provider,
		http.Cfg,
	)

	WireTestSet = wire.NewSet(
		CfgTest,
		Provider,
		ProviderSDKTestRepo,
		resolver.ValidatorsTest,
		dispatcher2.WireTestSet,
		plugins.WireTestSet,
		authentication.WireTestSet,
		http.Provider,
		http.CfgTest,
	)
)
