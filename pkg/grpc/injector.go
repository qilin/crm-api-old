// +build wireinject

package grpc

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
)

// Build
func Build(ctx context.Context, initial config.Initial, observer invoker.Observer) (*PoolManager, func(), error) {
	panic(wire.Build(provider.Set, WireSet, wire.Struct(new(provider.AwareSet), "*")))
}
