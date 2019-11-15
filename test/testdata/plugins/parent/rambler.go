package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path"
	"text/template"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
	"github.com/qilin/crm-api/internal/sdk/common"
	"github.com/qilin/crm-api/test/testdata/plugins/parent/rambler"
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

		// // get http request from context
		// req, ok := ctx.Value("request").(*http.Request)
		// if !ok {
		// 	log.Emergency("can't extract *http.Request from context")
		// }

		// // get cookie
		// cookie, err := req.Cookie(cfg["parent_auth_cookie_name"])
		// if err != nil {
		// 	return response, err
		// }

		// authToken := cookie.Value
		// // todo: verify authToken
		// if len(authToken) == 0 {
		// 	return response, errors.New("not authenticated")
		// }

		// issue JWT
		jwt, err := utils.IssueJWT("", "", "", request.QilinProductUUID, 0, keyPair.Private)
		if err != nil {
			return response, err
		}
		// url to return
		response.Meta = map[string]interface{}{
			"url": utils.AddURLParams(cfg["parent_iframe_url"], map[string]string{"jwt": string(jwt)}),
			// "cookie": authToken,
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
	r.GET("/integration/game/iframe", p.runTestGame)
	r.GET("/integration/game/billing", p.billingCallback)
	r.POST("/api/v2/svc/payment/create", p.createOrder)
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
		"GameUUID": ctx.QueryParam("uuid"),
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
