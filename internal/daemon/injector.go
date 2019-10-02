// +build wireinject

package daemon

import (
	"context"

	"github.com/google/wire"
	"github.com/qilin/go-core/config"
	"github.com/qilin/go-core/invoker"
	"github.com/qilin/go-core/provider"
)

// Build
func Build(ctx context.Context, initial config.Initial, observer invoker.Observer) (*Daemon, func(), error) {
	panic(wire.Build(WireSet, provider.Set, wire.Struct(new(provider.AwareSet), "*")))
}

// BuildTest
func BuildTest(ctx context.Context, initial config.Initial, observer invoker.Observer) (*Daemon, func(), error) {
	panic(wire.Build(WireTestSet, provider.Set, wire.Struct(new(provider.AwareSet), "*")))
}
