package dispatcher

import (
	"context"

	"github.com/qilin/crm-api/internal/auth"
	"github.com/qilin/crm-api/internal/jwt"

	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/qilin/crm-api/internal/dispatcher/common"
	"github.com/qilin/crm-api/pkg/graphql"
)

// Dispatcher
type Dispatcher struct {
	ctx     context.Context
	cfg     Config
	authCfg common.OAuth2
	appSet  AppSet
	provider.LMT
}

// dispatch
func (d *Dispatcher) Dispatch(echoHttp *echo.Echo) error {

	// middleware#2: recover
	echoHttp.Use(middleware.Recover())

	// middleware#1: CORS
	if d.cfg.Debug {
		echoHttp.Use(middleware.CORS())
	} else {
		echoHttp.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     d.cfg.CORS.Allowed,
			AllowMethods:     d.cfg.CORS.Methods,
			AllowHeaders:     d.cfg.CORS.Headers,
			AllowCredentials: false,
		}))
	}

	auth.New(d.LMT).RegisterAPIGroup(echoHttp)

	// init group routes
	grp := &common.Groups{
		Auth:    echoHttp.Group(common.AuthGroupPath),
		GraphQL: echoHttp.Group(common.GraphQLGroupPath),
		Common:  echoHttp,
		V1:      echoHttp.Group(common.V1Path),
	}

	d.graphqlGroup(grp.GraphQL)
	d.commonGroup(grp.Common)

	// init routes
	for _, handler := range d.appSet.Handlers {
		handler.Route(grp)
	}

	return nil
}

func (d *Dispatcher) graphqlGroup(grp *echo.Group) {
	// GraphQL JWT Middleware
	grp.Use(d.graphqlJWTMiddleware)
}

func (d *Dispatcher) commonGroup(grp *echo.Echo) {
	// add static or handlers
}

// Config
type Config struct {
	Debug   bool `fallback:"shared.debug"`
	WorkDir string
	OAuth   common.OAuth2
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
	GraphQL     *graphql.GraphQL
	Handlers    common.Handlers
	JwtVerifier *jwt.JWTVerefier
}

// New
func New(ctx context.Context, set provider.AwareSet, appSet AppSet, cfg *Config) *Dispatcher {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": common.Prefix})
	return &Dispatcher{
		ctx:    ctx,
		cfg:    *cfg,
		appSet: appSet,
		LMT:    &set,
	}
}
