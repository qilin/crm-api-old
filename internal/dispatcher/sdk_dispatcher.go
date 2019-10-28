package dispatcher

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/qilin/crm-api/internal/dispatcher/common"
)

// Dispatcher
type SDKDispatcher struct {
	ctx      context.Context
	cfg      Config
	handlers common.Handlers
	provider.LMT
}

// dispatch
func (d *SDKDispatcher) Dispatch(echoHttp *echo.Echo) error {

	echoHttp.Use(middleware.Logger())
	echoHttp.Use(middleware.Recover())

	echoHttp.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: d.cfg.CORS.Allowed,
		AllowMethods: d.cfg.CORS.Methods,
		AllowHeaders: d.cfg.CORS.Headers,
	}))

	// init group routes
	grp := &common.Groups{
		Common: echoHttp,
		V1:     echoHttp.Group(common.V1Path),
		SDK:    echoHttp.Group(common.SDKPath),
	}

	// init routes
	for _, handler := range d.handlers {
		handler.Route(grp)
	}

	return nil
}

// New
func NewSDK(ctx context.Context, set provider.AwareSet, h common.Handlers, cfg *Config) *SDKDispatcher {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": common.Prefix})
	return &SDKDispatcher{
		ctx:      ctx,
		cfg:      *cfg,
		handlers: h,
		LMT:      &set,
	}
}
