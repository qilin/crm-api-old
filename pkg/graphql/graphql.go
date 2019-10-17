package graphql

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/qilin/crm-api/internal/dispatcher/common"

	gqlgen "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/qilin/crm-api/internal/generated/graphql"
	gqErrs "github.com/qilin/crm-api/pkg/graphql/errors"
	"github.com/vektah/gqlparser/gqlerror"
)

var (
	ErrInternalServer = errors.New("internal server error")
	ErrAccessDenied   = errors.New("access denied")
)

// GraphQL
type GraphQL struct {
	ctx      context.Context
	resolver *graphql.Config
	cfg      *Config
	provider.LMT
}

// Route
func (g *GraphQL) Route(groups *common.Groups) {
	upgrader := websocket.Upgrader{}

	options := []handler.Option{
		handler.WebsocketUpgrader(upgrader),
		handler.IntrospectionEnabled(g.cfg.Introspection),
		handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
			if e, ok := err.(error); ok {
				return gqErrs.WrapPanicErr(e)
			}
			g.L().Alert("unhandled panic, err: %v", logger.Args(err))
			return nil
		}),
		handler.ErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
			switch e.(type) {
			case *gqErrs.PanicErr:
				g.L().Alert("recover on middleware, err: %v", logger.Args(e))
			case *gqErrs.AccessDeniedErr:
				g.L().Info("internal server error, err: %v", logger.Args(e))
			case *gqErrs.ClientErr:
				g.L().Error("internal server error, err: %v", logger.Args(e))
			default:
				g.L().Error("internal server error, err: %v", logger.Args(e))
				e = ErrInternalServer
			}
			return gqlgen.DefaultErrorPresenter(ctx, e)
		}),
	}

	if g.cfg.Debug {
		groups.GraphQL.Any(g.cfg.Playground.Route, echo.WrapHandler(handler.Playground(g.cfg.Playground.Name, g.cfg.Playground.Endpoint)))
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
		options = append(options, handler.RequestMiddleware(func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
			startTime := time.Now()
			rc := gqlgen.GetRequestContext(ctx)
			resp := next(ctx)
			e := strings.ReplaceAll(rc.Errors.Error(), "\n", " ")
			g.L().Debug("\nVARS:\n%+v\nQUERY:\n%v\nRESPONSE:\n%v\nERROR:\n%v\n",
				logger.Args(rc.Variables, strings.TrimRight(rc.RawQuery, "\n"), string(resp), e),
				logger.WithFields(logger.Fields{
					"time": time.Since(startTime).String(),
				}),
			)
			return resp
		}))
	}

	h := handler.GraphQL(
		graphql.NewExecutableSchema(*g.resolver),
		options...,
	)

	groups.GraphQL.Any(g.cfg.Route, echo.WrapHandler(h))
}

type PlaygroundCfg struct {
	Route    string
	Name     string
	Endpoint string
}

// Config
type Config struct {
	Debug         bool `fallback:"shared.debug"`
	Introspection bool
	Playground    PlaygroundCfg
	Route         string
	invoker       *invoker.Invoker
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
func New(ctx context.Context, resolver graphql.Config, set provider.AwareSet, cfg *Config) *GraphQL {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": Prefix})
	return &GraphQL{
		ctx:      ctx,
		resolver: &resolver,
		cfg:      cfg,
		LMT:      &set,
	}
}
