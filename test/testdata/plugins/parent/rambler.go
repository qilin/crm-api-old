package main

import (
	"bytes"
	"context"
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

type RamblerUserInfo struct {
	Id        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Status    string `json:"status"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
}

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

		// todo: issue JWT

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
