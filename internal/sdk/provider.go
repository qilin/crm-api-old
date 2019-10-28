package sdk

import (
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/db/repo"
	"github.com/qilin/crm-api/internal/db/trx"
	"github.com/qilin/crm-api/internal/resolver"
	"github.com/qilin/crm-api/pkg/postgres"
)

type AppSet struct {
	Repo Repo
	Trx  *trx.Manager
}

var (
	ProviderRepo = wire.NewSet(
		repo.NewPlatformRepo,
		trx.NewTrxManager,
	)

	ProviderSDKRepo = wire.NewSet(
		ProviderRepo,
		wire.Struct(new(Repo), "*"),
		postgres.WireSet,
	)

	ProviderSDKTestRepo = wire.NewSet(
		ProviderRepo,
		wire.Struct(new(Repo), "*"),
		resolver.ValidatorsTest,
		postgres.WireTestSet,
	)

	WireSet = wire.NewSet(
		ProviderSDKRepo,
		resolver.ValidatorsProduction,
		wire.Struct(new(AppSet), "*"),
	)

	WireTestSet = wire.NewSet(
		ProviderSDKTestRepo,
		resolver.ValidatorsTest,
		wire.Struct(new(AppSet), "*"),
	)
)
