package resolver

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/db/repo"
	"github.com/qilin/crm-api/internal/db/trx"
	"github.com/qilin/crm-api/internal/generated/graphql"
	"github.com/qilin/crm-api/internal/validators"
	"github.com/qilin/crm-api/pkg/postgres"
	"gopkg.in/go-playground/validator.v9"
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
	Trx  *trx.Manager
}

// Provider
func Provider(ctx context.Context, set provider.AwareSet, appSet AppSet, cfg *Config, validate *validator.Validate) (graphql.Config, func(), error) {
	c := New(ctx, set, appSet, cfg, validate)
	return c, func() {}, nil
}

var (
	ProviderRepo = wire.NewSet(
		repo.NewUserRepo,
		repo.NewJwtKeysRepo,
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

	ValidatorsProduction = wire.NewSet(
		validators.Provider,
		validators.ProviderValidators,
	)
	ValidatorsTest = wire.NewSet(
		validators.Provider,
		validators.ProviderValidators,
	)

	WireSet = wire.NewSet(
		Provider,
		Cfg,
		ProviderRepoProduction,
		ValidatorsProduction,
		wire.Struct(new(AppSet), "*"),
	)
	WireTestSet = wire.NewSet(
		Provider,
		CfgTest,
		ProviderTestRepo,
		ValidatorsTest,
		wire.Struct(new(AppSet), "*"),
	)
)
