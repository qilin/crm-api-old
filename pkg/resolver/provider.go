package resolver

import (
	"context"

	"github.com/google/wire"
	"github.com/qilin/crm-api/generated/graphql"
	"github.com/qilin/crm-api/internal/db/repo"
	"github.com/qilin/crm-api/internal/db/trx"
	"github.com/qilin/crm-api/internal/jwt"
	"github.com/qilin/crm-api/internal/validators"
	"github.com/qilin/crm-api/pkg/postgres"
	"github.com/qilin/go-core/config"
	"github.com/qilin/go-core/provider"
	validator "gopkg.in/go-playground/validator.v9"
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
func Provider(ctx context.Context, set provider.AwareSet, appSet AppSet, cfg *Config, validate *validator.Validate, jwt *jwt.Config) (graphql.Config, func(), error) {
	c := New(ctx, set, appSet, cfg, validate, jwt)
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
		validators.Provider,
		validators.ProviderValidators,
		jwt.Provider,
		Provider,
		Cfg,
		ProviderRepoProduction,
		wire.Struct(new(AppSet), "*"),
	)
	WireTestSet = wire.NewSet(
		validators.Provider,
		validators.ProviderValidators,
		jwt.Provider,
		Provider,
		CfgTest,
		ProviderTestRepo,
		wire.Struct(new(AppSet), "*"),
	)
)
