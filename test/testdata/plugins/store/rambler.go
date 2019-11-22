package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"text/template"

	"github.com/spf13/viper"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
	"github.com/qilin/crm-api/internal/sdk/common"
	"github.com/qilin/crm-api/pkg/qilin"
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
	PublicKey  string
	PrivateKey string
}

type RamblerID struct {
	Kid           string
	RsaPublicKey  string
	RsaPrivateKey string
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
	Qilin  string
}

type storeConfig struct {
	Auth   Auth
	Keys   Keys
	Routes Routes
	URL    URL
}

const (
	PluginName                    = "store"
	storeIndexTpl                 = "./web/store/store.html"
	storeIframeProviderTpl        = "./web/store/game.html"
	storeIframeProviderTplSandbox = "./web/store/game-sandbox.html"
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

	var err error
	jwtKeyPair, err = utils.DecodePemECDSA(p.config.Keys.JWT.PrivateKey, p.config.Keys.JWT.PublicKey)
	if err != nil {
		log.Emergency("plugin: can not parse key pair")
	}

	ramblerKeyPair, err = utils.DecodePemRSA(p.config.Keys.RamblerID.RsaPrivateKey, p.config.Keys.RamblerID.RsaPublicKey)
	if err != nil {
		log.Emergency("plugin: can not parse rsa key pair")
	}
}

func (p *plugin) Name() string {
	return PluginName
}

func (p *plugin) Auth(authenticate common.Authenticate) common.Authenticate {
	return func(ctx context.Context, request common.AuthRequest, token *jwt.Claims, log logger.Logger) (response common.AuthResponse, err error) {
		// fake auth
		if p.config.Auth.Fake {
			// issue JWT
			jwt, err := utils.IssueJWT("", "", "123", request.QilinProductUUID, 0, jwtKeyPair.Private)
			if err != nil {
				return response, err
			}
			// url to return
			return common.AuthResponse{
				Meta: map[string]interface{}{
					"url": utils.AddURLParams(qilin.IframeURL(p.config.URL.Qilin), map[string]string{"jwt": string(jwt)}),
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
			"url": utils.AddURLParams(qilin.IframeURL(p.config.URL.Qilin), map[string]string{"jwt": string(jwt)}),
			// "cookie": authToken,
		}
		return
	}
}

func (p *plugin) Http(ctx context.Context, r *echo.Echo, log logger.Logger) {
	r.GET("/store", func(c echo.Context) error {
		return p.IndexHandler(c, log)
	})
	r.GET("/game", func(c echo.Context) error {
		return p.IframeProviderHandler(c, log)
	})
	r.GET("/integration/game/iframe", p.runTestGame)
	r.GET("/integration/game/sdk/v1/auth", func(c echo.Context) error {
		jwt, err := utils.IssueJWT("", "", "123", "fa14b399-ae9b-4111-9c7f-0f1fe2cc1eb7", 0, jwtKeyPair.Private)
		if err != nil {
			return err
		}
		// url to return
		return c.JSON(http.StatusOK, &common.AuthResponse{
			Meta: map[string]interface{}{
				"url": utils.AddURLParams(qilin.IframeURL(p.config.URL.Qilin), map[string]string{"jwt": string(jwt)}),
			},
		})
	})
	r.GET("/integration/game/billing", p.billingCallback)
	r.POST("/api/v2/svc/payment/create", p.createOrder)
	r.POST("/confirmPayment", func(c echo.Context) error {
		return p.confirmPayment(c, p.config.URL.Qilin)
	})
	r.GET("/items", func(c echo.Context) error {
		return p.getItem(c, p.config.URL.Qilin)
	})

	r.GET("/rsid", p.runTestRsid)
	r.GET("/rsidx", p.runTestRsidx)
}

func (p *plugin) IframeProviderHandler(ctx echo.Context, log logger.Logger) error {
	tplName := path.Base(storeIframeProviderTpl)
	tpl, err := template.New(tplName).ParseFiles(storeIframeProviderTpl)
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
	fmt.Println("run successfully verified", ctx.Request().RequestURI)

	tplName := path.Base(storeIframeProviderTplSandbox)
	tpl, err := template.New(tplName).ParseFiles(storeIframeProviderTplSandbox)
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	err = tpl.ExecuteTemplate(buf, tplName, map[string]interface{}{
		"GameUUID": "fa14b399-ae9b-4111-9c7f-0f1fe2cc1eb7",
	})
	if err != nil {
		return err
	}
	return ctx.HTML(http.StatusOK, buf.String())
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

		item, err := p.queryItem(p.config.URL.Qilin, "fa14b399-ae9b-4111-9c7f-0f1fe2cc1eb7", ctx.QueryParam("item"))
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"response": item,
			// map[string]interface{}{
			// 	"title":     "50 золотых монет",
			// 	"photo_url": "https://ihcdn3.ioimg.org/iov6live/images/payments/payment_new/payment_packs_images/small_diamond.png",
			// 	"price":     50.0,
			// },
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

func (p *plugin) confirmPayment(ctx echo.Context, entry string) error {

	var params = make(map[string]string)
	if err := ctx.Bind(&params); err != nil {
		return err
	}

	fmt.Println(params)

	req := common.OrderRequest{
		GameID: params["gameId"],
		UserID: "123",
		ItemID: params["itemId"],
	}
	data, err := json.Marshal(&req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	resp, err := http.Post(qilin.OrderURL(entry), "application/json;charset=utf-8", bytes.NewReader(data))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if resp.StatusCode != http.StatusOK {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "payment proceeding failed",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{})

}

func (p *plugin) getItem(ctx echo.Context, entry string) error {
	var gameId = ctx.QueryParam("game_id")
	var itemId = ctx.QueryParam("item_id")
	item, err := p.queryItem(entry, gameId, itemId)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, item)
}

type Item struct {
	Title     string  `json:"title"`
	Photo_url string  `json:"photo_url"`
	Price     float32 `json:"price"`
}

func (p *plugin) queryItem(entry, gameId, itemId string) (*Item, error) {
	u := fmt.Sprintf("%s?item_id=%s&game_id=%s", qilin.ItemsURL(entry), itemId, gameId)

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("provider error")
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var item Item
	err = json.Unmarshal(d, &item)
	return &item, err
}

func (p *plugin) runTestRsid(ctx echo.Context) error {
	id := utils.IDClient{
		KID: p.config.Keys.RamblerID.Kid,
		Key: ramblerKeyPair.Private,
	}
	rsid := ctx.QueryParam("rsid")
	info := id.RamblerIdGetProfileInfo(rsid, ctx.Request().RemoteAddr, ctx.Request().UserAgent())
	return ctx.HTML(http.StatusOK, info)
}

func (p *plugin) runTestRsidx(ctx echo.Context) error {
	id := utils.IDClient{
		KID: p.config.Keys.RamblerID.Kid,
		Key: ramblerKeyPair.Private,
	}
	rsidx := ctx.QueryParam("rsidx")
	info := id.RamblerIdGetProfileInfoX(rsidx, ctx.Request().RemoteAddr, ctx.Request().UserAgent())
	return ctx.HTML(http.StatusOK, info)
}

func (p *plugin) runTestMeta(ctx echo.Context) error {
	id := utils.IDClient{
		KID: p.config.Keys.RamblerID.Kid,
		Key: ramblerKeyPair.Private,
	}
	info := id.RamblerMeta(ctx.Request().RemoteAddr, ctx.Request().UserAgent())
	return ctx.HTML(http.StatusOK, info)
}
