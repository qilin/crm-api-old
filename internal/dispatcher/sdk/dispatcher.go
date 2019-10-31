package sdk

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/qilin/crm-api/internal/dispatcher/common"
)

// Dispatcher
type Dispatcher struct {
	ctx      context.Context
	cfg      Config
	handlers common.Handlers
	provider.LMT
}

// dispatch
func (d *Dispatcher) Dispatch(echoHttp *echo.Echo) error {

	echoHttp.Use(middleware.Logger())
	//echoHttp.Use(middleware.Recover())

	echoHttp.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: d.cfg.CORS.Allowed,
		AllowMethods: d.cfg.CORS.Methods,
		AllowHeaders: d.cfg.CORS.Headers,
	}))

	grp := &common.Groups{
		SDK: echoHttp.Group(common.SDKPath),
	}

	for _, handler := range d.handlers {
		handler.Route(grp)
	}

	return nil
}

type Config struct {
	CORS    common.CORS
	invoker *invoker.Invoker
}

// OnReload
func (c *Config) OnReload(callback func(ctx context.Context)) {
	c.invoker.OnReload(callback)
}

// Reload
func (c *Config) Reload(ctx context.Context) {
	c.invoker.Reload(ctx)
}

// New
func New(ctx context.Context, set provider.AwareSet, h common.Handlers, cfg *Config) *Dispatcher {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": common.Prefix})
	return &Dispatcher{
		ctx:      ctx,
		cfg:      *cfg,
		handlers: h,
		LMT:      &set,
	}
}
