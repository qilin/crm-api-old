package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"text/template"

	"github.com/qilin/crm-api/test/testdata/plugins/parent/utils"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
	"github.com/qilin/crm-api/internal/sdk/common"
)

type plugin struct {
	//
}

var (
	Plugin  plugin
	keyPair utils.KeyPair
)

type RamblerAuthRequest struct {
	GameId    string `json:"game_id" query:"game_id"`
	Slug      string `json:"slug" query:"slug"`
	Timestamp string `json:"timestamp" query:"timestamp"`
	UserId    string `json:"user_id" query:"user_id"`
	Sig       string `json:"sig" query:"sig"`
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

		appSecret := cfg["parent_app_secret"]

		iframeURLString, ok := cfg["parent_iframe_url"]
		if !ok {
			log.Emergency("plugin: can not find config value `parent_iframe_url`")
		}
		iframeURL, err := url.Parse(iframeURLString)

		// Check rambler auth & signature
		u, err := url.Parse(request.URL)
		q := u.Query()
		h := md5.New()
		params := fmt.Sprintf("game_id=%s&slug=%s&timestamp-%s&user_id=%s&%s",
			q.Get("game_id"), q.Get("slug"), q.Get("timestamp"), q.Get("user_id"), appSecret)
		io.WriteString(h, params)
		sig := string(h.Sum(nil))
		if sig != q.Get("sig") {
			return common.AuthResponse{}, errors.New("incorrect signature")
		}

		// QilinProductUUID: "3d4ff5f9-8614-4524-ba4b-378a9fdb4594"

		// Issue JWT and add it to game's URL
		jwt, err := utils.IssueJWT("", "", q.Get("user_id"), "3d4ff5f9-8614-4524-ba4b-378a9fdb4594", 0, keyPair.Private)
		if err != nil {
			return common.AuthResponse{}, err
		}
		ifq := iframeURL.Query()
		ifq.Add("jwt", string(jwt))
		iframeURL.RawQuery = ifq.Encode()

		return common.AuthResponse{
			Meta: map[string]interface{}{
				"url": iframeURL.String(),
			},
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
