package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/qilin/crm-api/internal/sdk/common"
)

type plugin struct {
	//
}

var (
	Plugin plugin
)

func (p *plugin) Name() string {
	return "rgames.plugin"
}

func (p *plugin) Auth(authenticate common.Authenticate) common.Authenticate {
	return func(ctx context.Context, request common.AuthRequest, token *jwt.Claims, log logger.Logger) (common.AuthResponse, error) {
		url, err := p.getGameUrl()
		if err != nil {
			return common.AuthResponse{}, err
		}

		return common.AuthResponse{
			Meta: map[string]string{
				"mode": "rgames",
				"url":  url,
			},
		}, nil
	}
}

func (p *plugin) getGameUrl() (string, error) {
	r, err := http.NewRequest(http.MethodGet, "https://games.rambler.ru/api/v2/svc/games/1069/run", nil)
	if err != nil {
		return "", err
	}
	r.Header.Add("Cookie", "rsid=eyJleHRyYSI6eyJkYXRhIjoiTmp6RHlGODdDYTJMMWhCUzA2RUF6U2hNSUFNbXhvZ09WSncwbEhYTGdlVjlOZkJFSUR3QUVkNmgySzh5WGNkdGdCREpnWW16eHdOc3gxTHZvcmQyekFHdE1BZWlnaXZNTWlTRmJUNHpGMTR6QUsyZjJ0YXFyMllsamJmVE5md1Jtako5M3l3c0NnNElcLzRuRlFLaDNCMyt5d0NcL05nc2FWMFQ5OEtsYUJQTUUzYWM1Yk9cL2VoTDJSN2xNRG5xejVoZWFOZFZsMlZJc0toaTgxbmIrRmhXSWtmcXlCXC9zTW1xXC9UTT0iLCJlbmNfa2V5Ijoia2V5MSJ9LCJyc2lkIjoiYjQzOWMwMmM1NDNmMGMzMDQ3NmJmNzVlZGQwZjZkZTIifQ.v2.x")

	c := &http.Client{}
	resp, err := c.Do(r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var v = make(map[string]interface{})
	if err := json.Unmarshal(data, &v); err != nil {
		return "", err
	}

	return v["execution_endpoint"].(string), nil
}

func (p *plugin) Http(ctx context.Context, r *echo.Echo, log logger.Logger) {
	cfg, ok := ctx.Value("config").(map[string]string)
	if !ok {
		log.Error("plugin: can not cast context config to map[string]string")
	}
	_ = cfg

	r.GET("/rgames/sdk/v1/iframe", func(ctx echo.Context) error {
		data, err := ioutil.ReadFile("web/rgames/index.html")
		if err != nil {
			fmt.Println(err)
			return ctx.HTML(http.StatusInternalServerError, err.Error())
		}
		return ctx.HTML(http.StatusOK, string(data))
	})
	r.POST("/rgames/sdk/v1/sdk/v1/auth", func(ctx echo.Context) error {
		url, err := p.getGameUrl()
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, common.AuthResponse{
			Meta: map[string]string{
				"mode": "rgames-own",
				"url":  url,
			},
		})
	})

}
