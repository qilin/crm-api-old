// +build wireinject

package http

import (
	"context"

	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/dispatcher"
	"github.com/qilin/go-core/config"
	"github.com/qilin/go-core/invoker"
	"github.com/qilin/go-core/provider"
)

// Build
func Build(ctx context.Context, initial config.Initial, observer invoker.Observer) (*HTTP, func(), error) {
	panic(wire.Build(
		provider.Set,
		wire.Bind(new(Dispatcher), new(*dispatcher.Dispatcher)),
		WireSet,
		wire.Struct(new(provider.AwareSet), "*")),
	)
}

// BuildTest
func BuildTest(ctx context.Context, initial config.Initial, observer invoker.Observer) (*HTTP, func(), error) {
	panic(wire.Build(
		provider.Set,
		wire.Bind(new(Dispatcher), new(*dispatcher.Dispatcher)),
		WireTestSet,
		wire.Struct(new(provider.AwareSet), "*")),
	)
}
