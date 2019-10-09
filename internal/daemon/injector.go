// +build wireinject

package daemon

import (
	"context"
	"github.com/qilin/crm-api/pkg/http"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/dispatcher"
)

// Build
func BuildHTTP(ctx context.Context, initial config.Initial, observer invoker.Observer) (*http.HTTP, func(), error) {
	panic(wire.Build(
		provider.Set,
		wire.Bind(new(http.Dispatcher), new(*dispatcher.Dispatcher)),
		dispatcher.WireSet,
		wire.Struct(new(provider.AwareSet), "*"),
		http.WireSet,
	))
}

// BuildTest
func BuildHTTPTest(ctx context.Context, initial config.Initial, observer invoker.Observer) (*http.HTTP, func(), error) {
	panic(wire.Build(
		provider.Set,
		wire.Bind(new(http.Dispatcher), new(*dispatcher.Dispatcher)),
		dispatcher.WireTestSet,
		wire.Struct(new(provider.AwareSet), "*"),
		http.WireTestSet,
	))
}
