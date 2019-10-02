// +build wireinject

package stan

import (
	"github.com/google/wire"
	"github.com/qilin/go-core/entrypoint"
	"github.com/qilin/go-core/provider"
)

// Build
func Build(slave entrypoint.Slaver) (*Stan, func(), error) {
	panic(wire.Build(WireSet, provider.Set, wire.Struct(new(provider.AwareSet), "*")))
}

// BuildTest
func BuildTest(slave entrypoint.Slaver) (*Stan, func(), error) {
	panic(wire.Build(WireTestSet, provider.Set, wire.Struct(new(provider.AwareSet), "*")))
}
