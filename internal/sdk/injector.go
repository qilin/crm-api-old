// +build wireinject

package sdk

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/dispatcher"
	"github.com/qilin/crm-api/pkg/http"
)

// Build
func BuildHTTP(ctx context.Context, initial config.Initial, observer invoker.Observer) (*http.HTTP, func(), error) {
	panic(wire.Build(
		provider.Set,
		wire.Bind(new(http.Dispatcher), new(*dispatcher.SDKDispatcher)),
		wire.Struct(new(provider.AwareSet), "*"),
		http.SDKWireSet,
		WireSet,
	))
}

// BuildTest
func BuildHTTPTest(ctx context.Context, initial config.Initial, observer invoker.Observer) (*http.HTTP, func(), error) {
	panic(wire.Build(
		provider.Set,
		wire.Bind(new(http.Dispatcher), new(*dispatcher.SDKDispatcher)),
		wire.Struct(new(provider.AwareSet), "*"),
		http.SDKWireTestSet,
		WireTestSet,
	))
}
