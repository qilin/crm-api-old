// +build wireinject

package sdk

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/qilin/crm-api/internal/sdk/common"
	dispatcher2 "github.com/qilin/crm-api/internal/sdk/dispatcher"
	"github.com/qilin/crm-api/internal/sdk/handlers"
	sdk "github.com/qilin/crm-api/internal/sdk/sdk"
	"github.com/qilin/crm-api/pkg/http"
)

// Build
func BuildHTTP(ctx context.Context, initial config.Initial, observer invoker.Observer) (*http.HTTP, func(), error) {
	panic(wire.Build(
		provider.Set,
		wire.Bind(new(http.Dispatcher), new(*dispatcher2.Dispatcher)),
		wire.Bind(new(common.SDK), new(*sdk.SDK)),
		wire.Struct(new(provider.AwareSet), "*"),
		handlers.ProviderSDKHandlers,
		sdk.WireSet,
	))
}

// BuildTest
func BuildHTTPTest(ctx context.Context, initial config.Initial, observer invoker.Observer) (*http.HTTP, func(), error) {
	panic(wire.Build(
		provider.Set,
		wire.Bind(new(http.Dispatcher), new(*dispatcher2.Dispatcher)),
		wire.Bind(new(common.SDK), new(*sdk.SDK)),
		wire.Struct(new(provider.AwareSet), "*"),
		handlers.ProviderSDKHandlers,
		sdk.WireTestSet,
	))
}
