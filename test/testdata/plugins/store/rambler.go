package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"text/template"

	"github.com/spf13/viper"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
	"github.com/qilin/crm-api/internal/sdk/common"
	"github.com/qilin/crm-api/test/testdata/plugins/store/rambler"
	"github.com/qilin/crm-api/test/testdata/plugins/store/utils"
)

type plugin struct {
	config storeConfig
}

type Auth struct {
	Fake           bool
	CookieName     string
	RsidCookieName string
}

type JWT struct {
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
}

type RamblerID struct {
	Kid           string
	RsaPublicKey  *rsa.PublicKey
	RsaPrivateKey *rsa.PrivateKey
}

type Keys struct {
	JWT       JWT
	RamblerID RamblerID
}

type Routes struct {
	Index  string
	Iframe string
}

type URL struct {
	Iframe string
}

type storeConfig struct {
	Auth   Auth
	Keys   Keys
	Routes Routes
	URL    URL
}

const (
	PluginName             = "store"
	storeIndexTpl          = "./web/store/store.html"
	storeIframeProviderTpl = "./web/store/iframe.html"
)

var (
	Plugin         plugin
	jwtKeyPair     utils.KeyPair
	ramblerKeyPair utils.RSAKeyPair
)

func (p *plugin) Init(ctx context.Context, cfg *viper.Viper, log logger.Logger) {
	config := storeConfig{}
	cfg.UnmarshalKey(PluginName, &config)
	p.config = config
}

func (p *plugin) Name() string {
	return PluginName
}

func (p *plugin) Auth(authenticate common.Authenticate) common.Authenticate {
	return func(ctx context.Context, request common.AuthRequest, token *jwt.Claims, log logger.Logger) (response common.AuthResponse, err error) {
		// fake auth
		if p.config.Auth.Fake {
			jwt, err := utils.IssueJWT("", "", "123", request.QilinProductUUID, 0, jwtKeyPair.Private)
			if err != nil {
				return response, err
			}
			return common.AuthResponse{
				Meta: map[string]interface{}{
					"url": utils.AddURLParams(p.config.URL.Iframe, map[string]string{"jwt": string(jwt)}),
				},
			}, nil
		}

		// real auth

		// get rsid from meta
		//meta, ok := request.Meta.(map[string]interface{})
		//if !ok {
		//	return common.AuthResponse{}, errors.New("bad request. request.meta must be map[string]interface{}")
		//}
		//rsid, ok := meta["rsid"].(string)
		//if !ok {
		//	return common.AuthResponse{}, errors.New("bad request. request.meta[rsid] is undefined")
		//}

		// get rsid from cookie
		req, ok := ctx.Value("request").(*http.Request)
		if !ok {
			log.Emergency("can't extract *http.Request from context")
		}
		rsidCookie, err := req.Cookie(p.config.Auth.RsidCookieName)
		if err != nil {
			return common.AuthResponse{}, errors.New("bad request. rsid cookie is undefined")
		}

		rsid := rsidCookie.Value

		id := utils.IDClient{
			KID: p.config.Keys.RamblerID.Kid,
			Key: ramblerKeyPair.Private,
		}
		info := id.RamblerIdGetProfileInfo(rsid, req.RemoteAddr, req.UserAgent())

		//todo: parse info
		log.Info(info)
		// todo: set auth cookie?, but we also have rsid cookie

		// issue JWT
		jwt, err := utils.IssueJWT("", "", "123", request.QilinProductUUID, 0, jwtKeyPair.Private)
		if err != nil {
			return response, err
		}
		// url to return
		response.Meta = map[string]interface{}{
			"url": utils.AddURLParams(p.config.URL.Iframe, map[string]string{"jwt": string(jwt)}),
			// "cookie": authToken,
		}
		return
	}
}

func (p *plugin) Http(ctx context.Context, r *echo.Echo, log logger.Logger) {
	r.GET(p.config.Routes.Index, func(c echo.Context) error {
		return p.IndexHandler(c, log)
	})
	r.GET(p.config.Routes.Iframe, func(c echo.Context) error {
		log.Info(p.config.Routes.Iframe)
		return p.IframeProviderHandler(c, log)
	})
	r.GET("/integration/game/iframe", p.runTestGame)
	r.GET("/integration/game/billing", p.billingCallback)
	r.POST("/api/v2/svc/payment/create", p.createOrder)
}

func (p *plugin) IframeProviderHandler(ctx echo.Context, log logger.Logger) error {
	dir, _ := os.Getwd()
	log.Info(dir)
	log.Info(storeIframeProviderTpl)
	tplPath := storeIframeProviderTpl
	d, _ := ioutil.ReadDir("./web/store/")
	for _, f := range d {
		log.Info(f.Name())
	}

	tplName := path.Base(tplPath)
	tpl, err := template.New(tplName).ParseFiles(tplPath)
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	err = tpl.ExecuteTemplate(buf, tplName, map[string]interface{}{
		"GameUUID": ctx.QueryParam("uuid"),
	})
	if err != nil {
		return err
	}
	return ctx.HTML(http.StatusOK, buf.String())
}

func (p *plugin) IndexHandler(ctx echo.Context, log logger.Logger) error {
	tplPath := storeIndexTpl

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

func (p *plugin) runTestGame(ctx echo.Context) error {
	if !rambler.VerifySignature(ctx.QueryParams(), "6f12ff821d49e386c0918415322d0b74",
		"user_id", "game_id", "slug", "timestamp") {
		fmt.Println("bad signature", ctx.Request().RequestURI)
		return ctx.HTML(http.StatusUnauthorized, "Wrong Signature")
	}
	fmt.Println("billing callback successfully verified", ctx.Request().RequestURI)
	return ctx.HTML(http.StatusOK, `
<script src="//sandbox.games.rambler.ru/assets/ext/rgames.js" ></script>
<script>
rgames.init().then(() => {
	rgames.showOrderBox( {
		item : 100500 ,
		type : '' ,
	} ) ;
} ) ;	
</script>
		`)
}

func (p *plugin) billingCallback(ctx echo.Context) error {
	switch ctx.QueryParam("notification_type") {
	case "get_item":
		if !rambler.VerifySignature(ctx.QueryParams(), "6f12ff821d49e386c0918415322d0b74",
			"item", "app_id", "user_id", "receiver_id", "lang", "notification_type") {
			fmt.Println("bad signature", ctx.Request().RequestURI)
			return ctx.HTML(http.StatusUnauthorized, "Wrong Signature")
		}
		fmt.Println("billing callback successfully verified", ctx.Request().RequestURI)
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"response": map[string]interface{}{
				"title":     "50 золотых монет",
				"photo_url": "https://ihcdn3.ioimg.org/iov6live/images/payments/payment_new/payment_packs_images/small_diamond.png",
				"price":     50.0,
			},
		})
	case "order_status_change":
		if !rambler.VerifySignature(ctx.QueryParams(), "6f12ff821d49e386c0918415322d0b74",
			"item", "app_id", "user_id", "receiver_id", "lang", "order_id", "item_price", "status", "notification_type") {
			fmt.Println("bad signature", ctx.Request().RequestURI)
			return ctx.HTML(http.StatusUnauthorized, "Wrong Signature")
		}
		fmt.Println("billing callback successfully verified", ctx.Request().RequestURI)
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"response": map[string]interface{}{
				"order_id": ctx.QueryParam("order_id"),
			},
		})
	}
	return ctx.HTML(http.StatusNotFound, "method not found")
}

func (p *plugin) createOrder(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{})
}
