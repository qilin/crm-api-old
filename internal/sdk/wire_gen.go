// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package sdk

import (
	"context"
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/metric"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/ProtocolONE/go-core/v2/pkg/tracing"
	"github.com/qilin/crm-api/internal/authentication"
	"github.com/qilin/crm-api/internal/authentication/providers"
	"github.com/qilin/crm-api/internal/db/repo"
	"github.com/qilin/crm-api/internal/plugins"
	"github.com/qilin/crm-api/internal/sdk/dispatcher"
	"github.com/qilin/crm-api/internal/sdk/handlers"
	repo2 "github.com/qilin/crm-api/internal/sdk/repo"
	"github.com/qilin/crm-api/internal/sdk/sdk"
	"github.com/qilin/crm-api/internal/validators"
	"github.com/qilin/crm-api/pkg/http"
	"github.com/qilin/crm-api/pkg/postgres"
)

// Injectors from injector.go:

func BuildHTTP(ctx context.Context, initial config.Initial, observer invoker.Observer) (*http.HTTP, func(), error) {
	configurator, cleanup, err := config.Provider(initial, observer)
	if err != nil {
		return nil, nil, err
	}
	loggerConfig, cleanup2, err := logger.ProviderCfg(configurator)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	zap, cleanup3, err := logger.Provider(ctx, loggerConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	metricConfig, cleanup4, err := metric.ProviderCfg(configurator)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	scope, cleanup5, err := metric.ProviderPrometheus(ctx, zap, metricConfig)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	tracingConfig, cleanup6, err := tracing.ProviderCfg(configurator)
	if err != nil {
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	tracer, cleanup7, err := tracing.Provider(ctx, tracingConfig, zap)
	if err != nil {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	awareSet := provider.AwareSet{
		Logger: zap,
		Metric: scope,
		Tracer: tracer,
	}
	pluginsConfig, cleanup8, err := plugins.Cfg(configurator)
	if err != nil {
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	pluginManager, cleanup9, err := plugins.Provider(ctx, pluginsConfig, awareSet, initial)
	if err != nil {
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
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
		return nil, nil, err
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
		return nil, nil, err
	}
	authLogRepo := repo.NewAuthLogRepo(db)
	usersRepo := repo.NewUsersRepo(db)
	userProviderMapRepo := repo.NewAuthProviderRepo(db)
	appSet := authentication.AppSet{
		AuthLog:         authLogRepo,
		UsersRepo:       usersRepo,
		UserProviderMap: userProviderMapRepo,
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
		return nil, nil, err
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
		return nil, nil, err
	}
	providersConfig, cleanup14, err := providers.Cfg(configurator)
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
		return nil, nil, err
	}
	p1Provider, cleanup15, err := providers.ProviderP1(ctx, awareSet, providersConfig)
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
		return nil, nil, err
	}
	authenticationConfig, cleanup16, err := authentication.Cfg(configurator)
	if err != nil {
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
		return nil, nil, err
	}
	authenticationService, cleanup17, err := authentication.Provider(ctx, pluginManager, awareSet, appSet, validate, p1Provider, authenticationConfig)
	if err != nil {
		cleanup16()
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
		return nil, nil, err
	}
	storeRepo := repo.NewPlatformRepo(db)
	storeJWTKeyRepo := repo.NewPlatformJWTKeyRepo(db)
	storeGamesRepo := repo.NewProductsRepo(db)
	userMapRepo := repo.NewUserMapRepo(db)
	actionsLog := repo.ActionsLogProvider(db)
	repoRepo := &repo2.Repo{
		Store:       storeRepo,
		StoreJWTKey: storeJWTKeyRepo,
		StoreGames:  storeGamesRepo,
		UserMap:     userMapRepo,
		ActionsLog:  actionsLog,
	}
	sdkConfig, cleanup18, err := sdk.Cfg(configurator)
	if err != nil {
		cleanup17()
		cleanup16()
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
		return nil, nil, err
	}
	sdkSDK, cleanup19, err := sdk.Provider(ctx, pluginManager, awareSet, repoRepo, sdkConfig)
	if err != nil {
		cleanup18()
		cleanup17()
		cleanup16()
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
		return nil, nil, err
	}
	commonHandlers, cleanup20, err := handlers.ProviderSDKHandlers(validate, awareSet, sdkSDK)
	if err != nil {
		cleanup19()
		cleanup18()
		cleanup17()
		cleanup16()
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
		return nil, nil, err
	}
	dispatcherAppSet := dispatcher.AppSet{
		Authentication: authenticationService,
		Handlers:       commonHandlers,
	}
	dispatcherConfig, cleanup21, err := dispatcher.Cfg(configurator)
	if err != nil {
		cleanup20()
		cleanup19()
		cleanup18()
		cleanup17()
		cleanup16()
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
		return nil, nil, err
	}
	dispatcherDispatcher, cleanup22, err := dispatcher.ProviderDispatcher(ctx, awareSet, dispatcherAppSet, dispatcherConfig)
	if err != nil {
		cleanup21()
		cleanup20()
		cleanup19()
		cleanup18()
		cleanup17()
		cleanup16()
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
		return nil, nil, err
	}
	httpConfig, cleanup23, err := http.Cfg(configurator)
	if err != nil {
		cleanup22()
		cleanup21()
		cleanup20()
		cleanup19()
		cleanup18()
		cleanup17()
		cleanup16()
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
		return nil, nil, err
	}
	httpHTTP, cleanup24, err := http.Provider(ctx, awareSet, dispatcherDispatcher, httpConfig)
	if err != nil {
		cleanup23()
		cleanup22()
		cleanup21()
		cleanup20()
		cleanup19()
		cleanup18()
		cleanup17()
		cleanup16()
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
		return nil, nil, err
	}
	return httpHTTP, func() {
		cleanup24()
		cleanup23()
		cleanup22()
		cleanup21()
		cleanup20()
		cleanup19()
		cleanup18()
		cleanup17()
		cleanup16()
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

func BuildHTTPTest(ctx context.Context, initial config.Initial, observer invoker.Observer) (*http.HTTP, func(), error) {
	configurator, cleanup, err := config.Provider(initial, observer)
	if err != nil {
		return nil, nil, err
	}
	loggerConfig, cleanup2, err := logger.ProviderCfg(configurator)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	zap, cleanup3, err := logger.Provider(ctx, loggerConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	metricConfig, cleanup4, err := metric.ProviderCfg(configurator)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	scope, cleanup5, err := metric.ProviderPrometheus(ctx, zap, metricConfig)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	tracingConfig, cleanup6, err := tracing.ProviderCfg(configurator)
	if err != nil {
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	tracer, cleanup7, err := tracing.Provider(ctx, tracingConfig, zap)
	if err != nil {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	awareSet := provider.AwareSet{
		Logger: zap,
		Metric: scope,
		Tracer: tracer,
	}
	pluginsConfig, cleanup8, err := plugins.CfgTest()
	if err != nil {
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	pluginManager, cleanup9, err := plugins.Provider(ctx, pluginsConfig, awareSet, initial)
	if err != nil {
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
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
		return nil, nil, err
	}
	authLogRepo := repo.NewAuthLogRepo(db)
	usersRepo := repo.NewUsersRepo(db)
	userProviderMapRepo := repo.NewAuthProviderRepo(db)
	appSet := authentication.AppSet{
		AuthLog:         authLogRepo,
		UsersRepo:       usersRepo,
		UserProviderMap: userProviderMapRepo,
	}
	validatorSet, cleanup11, err := validators.Provider()
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
		return nil, nil, err
	}
	validate, cleanup12, err := validators.ProviderValidators(validatorSet)
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
		return nil, nil, err
	}
	providersConfig, cleanup13, err := providers.Cfg(configurator)
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
		return nil, nil, err
	}
	p1Provider, cleanup14, err := providers.ProviderP1(ctx, awareSet, providersConfig)
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
		return nil, nil, err
	}
	authenticationConfig, cleanup15, err := authentication.CfgTest()
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
		return nil, nil, err
	}
	authenticationService, cleanup16, err := authentication.Provider(ctx, pluginManager, awareSet, appSet, validate, p1Provider, authenticationConfig)
	if err != nil {
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
		return nil, nil, err
	}
	storeRepo := repo.NewPlatformRepo(db)
	storeJWTKeyRepo := repo.NewPlatformJWTKeyRepo(db)
	storeGamesRepo := repo.NewProductsRepo(db)
	userMapRepo := repo.NewUserMapRepo(db)
	actionsLog := repo.ActionsLogProvider(db)
	repoRepo := &repo2.Repo{
		Store:       storeRepo,
		StoreJWTKey: storeJWTKeyRepo,
		StoreGames:  storeGamesRepo,
		UserMap:     userMapRepo,
		ActionsLog:  actionsLog,
	}
	sdkConfig, cleanup17, err := sdk.CfgTest()
	if err != nil {
		cleanup16()
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
		return nil, nil, err
	}
	sdkSDK, cleanup18, err := sdk.Provider(ctx, pluginManager, awareSet, repoRepo, sdkConfig)
	if err != nil {
		cleanup17()
		cleanup16()
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
		return nil, nil, err
	}
	commonHandlers, cleanup19, err := handlers.ProviderSDKHandlers(validate, awareSet, sdkSDK)
	if err != nil {
		cleanup18()
		cleanup17()
		cleanup16()
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
		return nil, nil, err
	}
	dispatcherAppSet := dispatcher.AppSet{
		Authentication: authenticationService,
		Handlers:       commonHandlers,
	}
	dispatcherConfig, cleanup20, err := dispatcher.CfgTest()
	if err != nil {
		cleanup19()
		cleanup18()
		cleanup17()
		cleanup16()
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
		return nil, nil, err
	}
	dispatcherDispatcher, cleanup21, err := dispatcher.ProviderDispatcher(ctx, awareSet, dispatcherAppSet, dispatcherConfig)
	if err != nil {
		cleanup20()
		cleanup19()
		cleanup18()
		cleanup17()
		cleanup16()
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
		return nil, nil, err
	}
	httpConfig, cleanup22, err := http.CfgTest()
	if err != nil {
		cleanup21()
		cleanup20()
		cleanup19()
		cleanup18()
		cleanup17()
		cleanup16()
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
		return nil, nil, err
	}
	httpHTTP, cleanup23, err := http.Provider(ctx, awareSet, dispatcherDispatcher, httpConfig)
	if err != nil {
		cleanup22()
		cleanup21()
		cleanup20()
		cleanup19()
		cleanup18()
		cleanup17()
		cleanup16()
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
		return nil, nil, err
	}
	return httpHTTP, func() {
		cleanup23()
		cleanup22()
		cleanup21()
		cleanup20()
		cleanup19()
		cleanup18()
		cleanup17()
		cleanup16()
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
