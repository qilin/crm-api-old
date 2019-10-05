// +build wireinject

package casbin

import (
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/casbin/casbin"
	"github.com/google/wire"
)

// Build returns casbin.Enforcer instance with resolved dependencies
func Build(initial config.Initial, observer invoker.Observer) (*casbin.Enforcer, func(), error) {
	panic(wire.Build(WireSet, config.WireSet))
}

// BuildTest returns stub/mock instance casbin.Enforcer with resolved dependencies
func BuildTest(observer invoker.Observer, model Model, policy Policy) (*casbin.Enforcer, func(), error) {
	panic(wire.Build(WireTestSet))
}
