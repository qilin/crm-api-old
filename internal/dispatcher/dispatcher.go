package dispatcher

import (
	"context"

	"github.com/qilin/crm-api/internal/auth"
	"github.com/qilin/crm-api/internal/db/domain"

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
	ctx    context.Context
	cfg    Config
	appSet AppSet
	provider.LMT
}

// dispatch
func (d *Dispatcher) Dispatch(echoHttp *echo.Echo) error {

	echoHttp.Use(middleware.Logger())
	echoHttp.Use(middleware.Recover())

	echoHttp.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     d.cfg.CORS.Allowed,
		AllowMethods:     d.cfg.CORS.Methods,
		AllowHeaders:     d.cfg.CORS.Headers,
		AllowCredentials: true,
	}))

	// authorization
	auth, err := auth.New(d.LMT, &d.cfg.Auth, d.appSet.Users)
	if err != nil {
		return err
	}
	auth.RegisterHandlers(echoHttp)

	// configure graphql group
	var gql = echoHttp.Group(common.GraphQLGroupPath)
	gql.Use(auth.Middleware)

	// init group routes
	grp := &common.Groups{
		GraphQL: gql,
		Common:  echoHttp,
		V1:      echoHttp.Group(common.V1Path),
		SDK:     echoHttp.Group(common.SDKPath),
	}

	d.graphqlGroup(grp)
	d.commonGroup(grp.Common)

	// init routes
	for _, handler := range d.appSet.Handlers {
		handler.Route(grp)
	}

	return nil
}

func (d *Dispatcher) graphqlGroup(group *common.Groups) {
	// add graphql handlers
}

func (d *Dispatcher) commonGroup(grp *echo.Echo) {
	// add static or handlers
}

// Config
type Config struct {
	Auth    auth.Config
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
	GraphQL  *graphql.GraphQL
	Handlers common.Handlers
	Users    domain.UserRepo
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
