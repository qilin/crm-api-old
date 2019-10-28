// +build wireinject

package daemon

import (
	"context"
	"github.com/qilin/crm-api/internal/auth"
	"github.com/qilin/crm-api/internal/handlers"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/dispatcher"
	"github.com/qilin/crm-api/internal/eventbus"
	"github.com/qilin/crm-api/internal/stan"
	"github.com/qilin/crm-api/internal/webhooks"
	"github.com/qilin/crm-api/pkg/http"
)

// Build
func BuildHTTP(ctx context.Context, initial config.Initial, observer invoker.Observer) (*http.HTTP, func(), error) {
	panic(wire.Build(
		provider.Set,
		wire.Bind(new(http.Dispatcher), new(*dispatcher.Dispatcher)),
		wire.Struct(new(provider.AwareSet), "*"),
		http.WireSet,
		stan.WireSet,
		eventbus.WireSet,
		webhooks.WireSet,
		wire.Struct(new(handlers.Handlers), "*"),
		handlers.ProviderHandlers,
		auth.WireSet,
		dispatcher.WireSet,
	))
}

// BuildTest
func BuildHTTPTest(ctx context.Context, initial config.Initial, observer invoker.Observer) (*http.HTTP, func(), error) {
	panic(wire.Build(
		provider.Set,
		wire.Bind(new(http.Dispatcher), new(*dispatcher.Dispatcher)),
		wire.Struct(new(provider.AwareSet), "*"),
		http.WireTestSet,
		stan.WireTestSet,
		eventbus.WireTestSet,
		webhooks.WireTestSet,
		wire.Struct(new(handlers.Handlers), "*"),
		handlers.ProviderHandlers,
		auth.WireSet,
		dispatcher.WireSet,
	))
}
