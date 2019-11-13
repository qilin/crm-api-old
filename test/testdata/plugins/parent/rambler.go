package main

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"path"
	"text/template"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
	"github.com/qilin/crm-api/internal/sdk/common"
	"github.com/qilin/crm-api/test/testdata/plugins/parent/utils"
)

type plugin struct {
	//
}

var (
	Plugin  plugin
	keyPair utils.KeyPair
)

func (p *plugin) Init(ctx context.Context, cfg map[string]string, log logger.Logger) {
	// load encryption keys
	sk, ok := cfg["parent_private_key"]
	if !ok {
		log.Emergency("plugin: can not load parent_private_key")
	}
	pk, ok := cfg["parent_public_key"]
	if !ok {
		log.Emergency("plugin: can not load parent_public_key")
	}

	var err error
	keyPair, err = utils.DecodePemECDSA(sk, pk)
	if err != nil {
		log.Emergency("plugin: can not parse key pair")
	}
}

func (p *plugin) Name() string {
	return "rambler.plugin"
}

func (p *plugin) Auth(authenticate common.Authenticate) common.Authenticate {
	return func(ctx context.Context, request common.AuthRequest, token *jwt.Claims, log logger.Logger) (response common.AuthResponse, err error) {
		cfg, ok := ctx.Value("config").(map[string]string)
		if !ok {
			log.Emergency("plugin: can not cast context config to map[string]string")
		}

		// fake auth
		if fake, ok := cfg["parent_auth_fake"]; ok && fake == "true" {
			url, ok := cfg["parent_iframe_url"]
			if !ok {
				url = "%parent_iframe_url%"
			}
			return common.AuthResponse{
				Meta: map[string]interface{}{
					"url": url,
				},
			}, nil
		}

		// get http request from context
		req, ok := ctx.Value("request").(*http.Request)
		if !ok {
			log.Emergency("can't extract *http.Request from context")
		}

		// get cookie
		cookie, err := req.Cookie(cfg["parent_auth_cookie_name"])
		if err != nil {
			return response, err
		}

		authToken := cookie.Value
		// todo: verify authToken
		if len(authToken) == 0 {
			return response, errors.New("not authenticated")
		}

		// issue JWT
		jwt, err := utils.IssueJWT("", "", "", "3d4ff5f9-8614-4524-ba4b-378a9fdb4594", 0, keyPair.Private)
		if err != nil {
			return response, err
		}
		// url to return
		response.Meta = map[string]interface{}{
			"url":    utils.AddURLParams(cfg["parent_iframe_url"], map[string]string{"jwt": string(jwt)}),
			"cookie": authToken,
		}
		return
	}
}

func (p *plugin) Http(ctx context.Context, r *echo.Echo, log logger.Logger) {
	cfg, ok := ctx.Value("config").(map[string]string)
	if !ok {
		log.Emergency("plugin: can not cast context config to map[string]string")
	}

	// Parent Iframe provider
	indexRoute, ok := cfg["parent_index_route"]
	if !ok {
		log.Emergency("plugin: can not find parent_index_route in config")
	}
	// Parent Iframe provider
	iframeProviderRoute, ok := cfg["parent_iframe_route"]
	if !ok {
		log.Emergency("plugin: can not find parent_iframe_route in config")
	}

	r.GET(indexRoute, func(c echo.Context) error {
		return p.IndexHandler(c, cfg, log)
	})
	r.GET(iframeProviderRoute, func(c echo.Context) error {
		return p.IframeProviderHandler(c, cfg, log)
	})
}

func (p *plugin) IframeProviderHandler(ctx echo.Context, cfg map[string]string, log logger.Logger) error {
	tplPath, ok := cfg["parent_iframe_template"]
	if !ok {
		log.Emergency("plugin: can not find parent_iframe_template in config")
	}

	tplName := path.Base(tplPath)
	tpl, err := template.New(tplName).ParseFiles(tplPath)
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	err = tpl.ExecuteTemplate(buf, tplName, map[string]interface{}{
		"QilinGameProxyURL": "",
		"IframeURL":         "",
	})
	if err != nil {
		return err
	}
	return ctx.HTML(http.StatusOK, buf.String())
}

func (p *plugin) IndexHandler(ctx echo.Context, cfg map[string]string, log logger.Logger) error {
	tplPath, ok := cfg["parent_index_template"]
	if !ok {
		log.Emergency("plugin: can not find parent_index_template in config")
	}

	tplName := path.Base(tplPath)
	tpl, err := template.New(tplName).ParseFiles(tplPath)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	buf := &bytes.Buffer{}
	err = tpl.ExecuteTemplate(buf, tplName, map[string]interface{}{})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return ctx.HTML(http.StatusOK, buf.String())
}
