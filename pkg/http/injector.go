// +build wireinject

package http

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/dispatcher"
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
