// +build wireinject

package graphql

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
)

// Build
func Build(ctx context.Context, initial config.Initial, observer invoker.Observer) (*GraphQL, func(), error) {
	panic(wire.Build(provider.Set, WireSet, wire.Struct(new(provider.AwareSet), "*")))
}

// BuildTest
func BuildTest(ctx context.Context, initial config.Initial, observer invoker.Observer) (*GraphQL, func(), error) {
	panic(wire.Build(provider.Set, WireTestSet, wire.Struct(new(provider.AwareSet), "*")))
}
