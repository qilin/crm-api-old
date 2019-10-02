// +build wireinject

package grpc

import (
	"context"

	"github.com/google/wire"
	"github.com/qilin/go-core/config"
	"github.com/qilin/go-core/invoker"
	"github.com/qilin/go-core/provider"
)

// Build
func Build(ctx context.Context, initial config.Initial, observer invoker.Observer) (*PoolManager, func(), error) {
	panic(wire.Build(provider.Set, WireSet, wire.Struct(new(provider.AwareSet), "*")))
}
