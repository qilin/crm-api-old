package plugins

import (
	"context"

	"github.com/spf13/viper"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/qilin/crm-api/internal/sdk/common"
)

type Httper interface {
	Name() string
	Http(ctx context.Context, r *echo.Echo, log logger.Logger)
}

type Authenticator interface {
	Provider() string
	Name() string
	Auth(authenticate common.Authenticate) common.Authenticate
}

type Orderer interface {
	Name() string
	Order(order common.Order) common.Order
}

type Initable interface {
	Name() string
	Init(ctx context.Context, cfg *viper.Viper, log logger.Logger)
}
