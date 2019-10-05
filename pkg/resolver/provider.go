package resolver

import (
	"context"

	"github.com/qilin/crm-api/internal/validators"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/generated/graphql"
	"github.com/qilin/crm-api/internal/db/repo"
	"github.com/qilin/crm-api/internal/db/trx"
	"github.com/qilin/crm-api/pkg/postgres"
)

// Cfg
func Cfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{}
	e := cfg.UnmarshalKeyOnReload(UnmarshalKey, c)
	return c, func() {}, e
}

// CfgTest
func CfgTest() (*Config, func(), error) {
	return &Config{}, func() {}, nil
}

type AppSet struct {
	Repo Repo
	//Validate validator.Validate
	Trx *trx.Manager
}

// Provider
func Provider(ctx context.Context, set provider.AwareSet, appSet AppSet, cfg *Config) (graphql.Config, func(), error) {
	c := New(ctx, set, appSet, cfg)
	return c, func() {}, nil
}

var (
	ProviderRepo = wire.NewSet(
		repo.NewListRepo,
		repo.NewUserRepo,
		trx.NewTrxManager,
	)
	ProviderRepoProduction = wire.NewSet(
		ProviderRepo,
		wire.Struct(new(Repo), "*"),
		postgres.WireSet,
	)
	ProviderTestRepo = wire.NewSet(
		ProviderRepo,
		wire.Struct(new(Repo), "*"),
		postgres.WireTestSet,
	)
	WireSet = wire.NewSet(
		Provider,
		Cfg,
		ProviderRepoProduction,
		wire.Struct(new(AppSet), "*"),
	)
	WireTestSet = wire.NewSet(
		Provider,
		CfgTest,
		ProviderTestRepo,
		validators.WireTestSet,
		wire.Struct(new(AppSet), "*"),
	)
)
