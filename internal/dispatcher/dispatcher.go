package dispatcher

import (
	"context"
	"fmt"
	"net/http"

	"github.com/qilin/crm-api/internal/auth"
	"github.com/qilin/crm-api/internal/jwt"

	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/qilin/crm-api/internal/dispatcher/common"
	"github.com/qilin/crm-api/pkg/graphql"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

	auth.New().RegisterAPIGroup(echoHttp)

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

	// init group routes
	grp := &common.Groups{
		Auth:    echoHttp.Group(common.AuthGroupPath),
		GraphQL: echoHttp.Group(common.GraphQLGroupPath),
		Common:  echoHttp,
		V1:      echoHttp.Group(common.V1Path),
	}

	d.graphqlGroup(grp.GraphQL)
	d.commonGroup(grp.Common)
	// auth.RegisterAuthGroup()/

	fmt.Println(d.cfg.Debug)

	fmt.Println("handlers", d.appSet.Handlers)
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

var googleCfg = &oauth2.Config{
	RedirectURL:  "http://localhost:8082/auth/oauth/callback",
	ClientID:     "633585228079-ssip47fknk77sfc2f930r71hjq010cra.apps.googleusercontent.com",
	ClientSecret: "SQmbV-Z3hmfY3UfDSlnEKGk-",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "openid"},
	Endpoint:     google.Endpoint,
}

var auth1Cfg = &oauth2.Config{
	RedirectURL:  "http://localhost:8082/auth/oauth/callback",
	ClientID:     "5da469eb2b13220001efe15e",
	ClientSecret: "mb1XDFcISlwwq7Vk2ybCJefrhiBZlpa6omCOhN9z4u3cMxupzQdP1N9PkWQzPS4A",
	Scopes:       []string{"openid"},
	Endpoint: oauth2.Endpoint{
		AuthURL:   "https://auth1.tst.protocol.one/oauth2/auth",
		TokenURL:  "https://auth1.tst.protocol.one/oauth2/token",
		AuthStyle: oauth2.AuthStyleInHeader,
	},
}

var openidProviders = make(map[string]*oidc.Provider)

func (d *Dispatcher) authGroup(grp *echo.Group) {

	grp.GET("/oauth/:provider", func(c echo.Context) error {
		var provider = c.Param("provider")
		switch provider {
		case "auth1":
			url := auth1Cfg.AuthCodeURL("some state XXX")
			return c.Redirect(http.StatusFound, url)
		case "google":
			url := googleCfg.AuthCodeURL("some state")
			return c.Redirect(http.StatusFound, url)
		}
		return fmt.Errorf("unknown provider '%s'", provider)
	})

	grp.GET("/oauth/callback", func(c echo.Context) error {
		var (
			state = c.FormValue("state")
			code  = c.FormValue("code")
		)
		fmt.Println(state)

		token, err := auth1Cfg.Exchange(oauth2.NoContext, code)
		if err != nil {
			return fmt.Errorf("code exchange failed: %s", err.Error())
		}

		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok {
			return fmt.Errorf("id_token not provided")
		}

		ctx := context.Background()
		provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
		if err != nil {
			// handle error
		}

		_ = rawIDToken
		_ = provider

		// Parse and verify ID Token payload.
		// idToken, err := verifier.Verify(ctx, rawIDToken)
		// if err != nil {
		// 	// handle error
		// }

		// // Extract custom claims
		// var claims struct {
		// 	Email    string `json:"email"`
		// 	Verified bool   `json:"email_verified"`
		// }
		// if err := idToken.Claims(&claims); err != nil {
		// 	// handle error
		// }
		return nil
	})

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
