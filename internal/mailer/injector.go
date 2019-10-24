// +build wireinject

package mailer

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/eventbus"
	"github.com/qilin/crm-api/internal/stan"
)

// Build
func BuildMailer(ctx context.Context, initial config.Initial, observer invoker.Observer) (*eventbus.EventBus, func(), error) {
	panic(wire.Build(
		provider.Set,
		wire.Struct(new(provider.AwareSet), "*"),
		stan.WireSet,
		eventbus.WireSet,
	))
}

// BuildTest
func BuildMailerTest(ctx context.Context, initial config.Initial, observer invoker.Observer) (*eventbus.EventBus, func(), error) {
	panic(wire.Build(
		provider.Set,
		wire.Struct(new(provider.AwareSet), "*"),
		stan.WireTestSet,
		eventbus.WireTestSet,
	))
}
