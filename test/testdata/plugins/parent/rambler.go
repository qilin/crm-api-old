package main

import (
	"context"
	"net/http"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
	"github.com/qilin/crm-api/internal/sdk/common"
)

type plugin struct {
	//
}

var (
	Plugin plugin
)

func (p *plugin) Name() string {
	return "example.plugin"
}

func (p *plugin) Auth(authenticate common.Authenticate) common.Authenticate {
	return func(ctx context.Context, request common.AuthRequest, token *jwt.Claims, log logger.Logger) (response common.AuthResponse, err error) {
		cfg, ok := ctx.Value("config").(map[string]string)
		if !ok {
			log.Emergency("plugin: can not cast context config to map[string]string")
		}

		url, ok := cfg["parent_iframe_url"]
		if !ok {
			url = "%parent_iframe_url%"
		}

		meta := map[string]interface{}{
			"url": url,
		}

		return common.AuthResponse{
			Meta: meta,
		}, nil
	}
}

func (p *plugin) Http(ctx context.Context, r *echo.Echo, log logger.Logger) {
	cfg, ok := ctx.Value("config").(map[string]string)
	if !ok {
		log.Emergency("plugin: can not cast context config to map[string]string")
	}

	route, ok := cfg["parent_http_route"]
	if !ok {
		log.Emergency("plugin: can not find parent_http_route in config")
	}

	// todo: add html page

	r.GET(route, func(c echo.Context) error {
		return c.JSON(http.StatusOK, "plugin endpoint")
	})
}
