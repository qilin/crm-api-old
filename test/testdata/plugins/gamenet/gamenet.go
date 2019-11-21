package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/qilin/crm-api/internal/sdk/common"
	"github.com/qilin/crm-api/test/testdata/plugins/store/rambler"
)

type plugin struct {
	//
}

var (
	Plugin plugin
)

func (p *plugin) Name() string {
	return "gamenet.plugin"
}

func (p *plugin) Auth(authenticate common.Authenticate) common.Authenticate {
	return func(ctx context.Context, request common.AuthRequest, token *jwt.Claims, log logger.Logger) (response common.AuthResponse, err error) {
		meta := map[string]string{
			"mode": "gamenet",
			"url":  "/games/khanwars/iframe?wmode=opaque",
		}

		return common.AuthResponse{
			Meta: meta,
		}, nil
	}
}

func (p *plugin) Http(ctx context.Context, r *echo.Echo, log logger.Logger) {
	cfg, ok := ctx.Value("config").(map[string]string)
	if !ok {
		log.Error("plugin: can not cast context config to map[string]string")
	}
	_ = cfg

	r.GET("/gamenet/sdk/v1/iframe", func(ctx echo.Context) error {
		// TODO user from jws token
		u, err := rambler.SignUrl(
			"https://gameplatform.stg.gamenet.ru/iframe/1095/qilin/?game_id=1&slug=protocolone&timestamp=1573646629488&user_id=1683086",
			"msBzPSaWYQ0piSLZJJNg")
		if err != nil {
			log.Error(err.Error())
			return ctx.HTML(http.StatusInternalServerError, "Internal Server Error")
		}

		return ctx.Redirect(http.StatusTemporaryRedirect, u)
	})
	r.GET("/gamenet/sdk/v1/items", func(ctx echo.Context) error {
		var gameId = ctx.QueryParam("game_id")
		var itemId = ctx.QueryParam("item_id")
		_ = gameId
		u, err := rambler.SignUrl(fmt.Sprintf(
			"https://api.stg.gamenet.ru/?method=gameplatform.qilin&item=%s&app_id=1&user_id=1683086&receiver_id=1683086&lang=ru_RU&notification_type=get_item",
			itemId,
		), "msBzPSaWYQ0piSLZJJNg")
		if err != nil {
			log.Error(err.Error())
			return ctx.HTML(http.StatusInternalServerError, "Internal Server Error")
		}

		resp, err := http.Get(u)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("provider error")
		}

		d, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var v = make(map[string]interface{})
		if err := json.Unmarshal(d, &v); err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, v["response"])
	})

	r.POST("/gamenet/sdk/v1/order", func(ctx echo.Context) error {
		var params = make(map[string]string)
		if err := ctx.Bind(&params); err != nil {
			log.Error(err.Error())
			return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "bad request",
			})
		}

		fmt.Println(params)
		u, err := rambler.SignUrl(fmt.Sprintf(
			"https://api.stg.gamenet.ru/?method=gameplatform.qilin&notification_type=order_status_change&item=%s&app_id=1&user_id=1683086&receiver_id=1683086&lang=ru_RU&order_id=%d&item_price=115.0&status=chargeable",
			params["item_id"],
			time.Now().Unix()-1574170567, // order_id
		), "msBzPSaWYQ0piSLZJJNg")

		if err != nil {
			log.Error(err.Error())
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "internal error",
			})
		}

		fmt.Println(u)
		resp, err := http.Get(u)
		if err != nil {
			log.Error(err.Error())
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "internal error",
			})
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Println(resp.Status)
			data, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(data))
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "payment failed",
			})
		}

		return ctx.JSON(http.StatusOK, map[string]interface{}{})
	})
}
