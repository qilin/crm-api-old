// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package resolver

import (
	"context"
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/metric"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/ProtocolONE/go-core/v2/pkg/tracing"
	"github.com/qilin/crm-api/internal/db/repo"
	"github.com/qilin/crm-api/internal/db/trx"
	"github.com/qilin/crm-api/internal/generated/graphql"
	"github.com/qilin/crm-api/internal/validators"
	"github.com/qilin/crm-api/pkg/cache"
	"github.com/qilin/crm-api/pkg/postgres"
)

// Injectors from injector.go:

func Build(ctx context.Context, initial config.Initial, observer invoker.Observer) (graphql.Config, func(), error) {
	configurator, cleanup, err := config.Provider(initial, observer)
	if err != nil {
		return graphql.Config{}, nil, err
	}
	loggerConfig, cleanup2, err := logger.ProviderCfg(configurator)
	if err != nil {
		cleanup()
		return graphql.Config{}, nil, err
	}
	zap, cleanup3, err := logger.Provider(ctx, loggerConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	metricConfig, cleanup4, err := metric.ProviderCfg(configurator)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	scope, cleanup5, err := metric.ProviderPrometheus(ctx, zap, metricConfig)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	tracingConfig, cleanup6, err := tracing.ProviderCfg(configurator)
	if err != nil {
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	tracer, cleanup7, err := tracing.Provider(ctx, tracingConfig, zap)
	if err != nil {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	awareSet := provider.AwareSet{
		Logger: zap,
		Metric: scope,
		Tracer: tracer,
	}
	cacheConfig, cleanup8, err := cache.ProviderCfg(configurator)
	if err != nil {
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	cacheCache, cleanup9, err := cache.Provider(ctx, awareSet, cacheConfig)
	if err != nil {
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	postgresConfig, cleanup10, err := postgres.ProviderCfg(configurator)
	if err != nil {
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	db, cleanup11, err := postgres.ProviderGORM(ctx, zap, postgresConfig)
	if err != nil {
		cleanup10()
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	userRepo := repo.NewUserRepo(db)
	usersRepo := repo.NewUsersRepo(db)
	gamesRepo := repo.NewGamesRepo(db)
	storefrontRepo := repo.NewStorefrontRepo(db)
	resolverRepo := Repo{
		User:        userRepo,
		Users:       usersRepo,
		Games:       gamesRepo,
		Storefronts: storefrontRepo,
	}
	manager := trx.NewTrxManager(db)
	appSet := AppSet{
		Cache: cacheCache,
		Repo:  resolverRepo,
		Trx:   manager,
	}
	resolverConfig, cleanup12, err := Cfg(configurator)
	if err != nil {
		cleanup11()
		cleanup10()
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	validatorSet, cleanup13, err := validators.Provider()
	if err != nil {
		cleanup12()
		cleanup11()
		cleanup10()
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	validate, cleanup14, err := validators.ProviderValidators(validatorSet)
	if err != nil {
		cleanup13()
		cleanup12()
		cleanup11()
		cleanup10()
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	graphqlConfig, cleanup15, err := Provider(ctx, awareSet, appSet, resolverConfig, validate)
	if err != nil {
		cleanup14()
		cleanup13()
		cleanup12()
		cleanup11()
		cleanup10()
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	return graphqlConfig, func() {
		cleanup15()
		cleanup14()
		cleanup13()
		cleanup12()
		cleanup11()
		cleanup10()
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

func BuildTest(ctx context.Context, initial config.Initial, observer invoker.Observer) (graphql.Config, func(), error) {
	configurator, cleanup, err := config.Provider(initial, observer)
	if err != nil {
		return graphql.Config{}, nil, err
	}
	loggerConfig, cleanup2, err := logger.ProviderCfg(configurator)
	if err != nil {
		cleanup()
		return graphql.Config{}, nil, err
	}
	zap, cleanup3, err := logger.Provider(ctx, loggerConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	metricConfig, cleanup4, err := metric.ProviderCfg(configurator)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	scope, cleanup5, err := metric.ProviderPrometheus(ctx, zap, metricConfig)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	tracingConfig, cleanup6, err := tracing.ProviderCfg(configurator)
	if err != nil {
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	tracer, cleanup7, err := tracing.Provider(ctx, tracingConfig, zap)
	if err != nil {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	awareSet := provider.AwareSet{
		Logger: zap,
		Metric: scope,
		Tracer: tracer,
	}
	cacheConfig, cleanup8, err := cache.ProviderCfg(configurator)
	if err != nil {
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	cacheCache, cleanup9, err := cache.Provider(ctx, awareSet, cacheConfig)
	if err != nil {
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	db, cleanup10, err := postgres.ProviderGORMTest()
	if err != nil {
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	userRepo := repo.NewUserRepo(db)
	usersRepo := repo.NewUsersRepo(db)
	gamesRepo := repo.NewGamesRepo(db)
	storefrontRepo := repo.NewStorefrontRepo(db)
	resolverRepo := Repo{
		User:        userRepo,
		Users:       usersRepo,
		Games:       gamesRepo,
		Storefronts: storefrontRepo,
	}
	manager := trx.NewTrxManager(db)
	appSet := AppSet{
		Cache: cacheCache,
		Repo:  resolverRepo,
		Trx:   manager,
	}
	resolverConfig, cleanup11, err := CfgTest()
	if err != nil {
		cleanup10()
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	validatorSet, cleanup12, err := validators.Provider()
	if err != nil {
		cleanup11()
		cleanup10()
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	validate, cleanup13, err := validators.ProviderValidators(validatorSet)
	if err != nil {
		cleanup12()
		cleanup11()
		cleanup10()
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	graphqlConfig, cleanup14, err := Provider(ctx, awareSet, appSet, resolverConfig, validate)
	if err != nil {
		cleanup13()
		cleanup12()
		cleanup11()
		cleanup10()
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return graphql.Config{}, nil, err
	}
	return graphqlConfig, func() {
		cleanup14()
		cleanup13()
		cleanup12()
		cleanup11()
		cleanup10()
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
