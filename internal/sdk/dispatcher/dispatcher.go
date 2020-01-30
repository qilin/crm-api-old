package dispatcher

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/qilin/crm-api/internal/authentication"
	"github.com/qilin/crm-api/internal/dispatcher/common"
	common2 "github.com/qilin/crm-api/internal/handlers/common"
)

// Dispatcher
type Dispatcher struct {
	ctx context.Context
	cfg Config
	app AppSet
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

	// middleware: session to context
	echoHttp.Use(d.app.Authentication.Middleware)

	v1 := echoHttp.Group(common.V1Path)

	grp := &common2.Groups{
		Auth:   v1.Group(common.AuthGroupPath),
		SDK:    echoHttp.Group(common.SDKPath),
		Common: echoHttp,
	}

	// auth routes
	d.app.Authentication.Route(grp)

	for _, handler := range d.app.Handlers {
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

type AppSet struct {
	Authentication *authentication.AuthenticationService
	Handlers       common2.Handlers
}

// New
func New(ctx context.Context, set provider.AwareSet, appSet AppSet, cfg *Config) *Dispatcher {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": common.Prefix})
	return &Dispatcher{
		ctx: ctx,
		cfg: *cfg,
		app: appSet,
		LMT: &set,
	}
}
