package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/qilin/crm-api/internal/handlers/common"

	"github.com/qilin/crm-api/internal/db/domain/store"
	"github.com/qilin/crm-api/internal/db/repo"
)

type Internal struct {
	games       *repo.GamesRepo
	storefronts *repo.StorefrontRepo
}

func NewInternal(games *repo.GamesRepo, storefronts *repo.StorefrontRepo) *Internal {
	return &Internal{
		games:       games,
		storefronts: storefronts,
	}
}

func (h *Internal) Route(groups *common.Groups) {
	groups.V1.POST("/internal/games", h.publishGames)
	groups.V1.POST("/internal/modules", h.saveModule)
}

func (h *Internal) publishGames(ctx echo.Context) error {
	var games store.GameSlice
	if err := ctx.Bind(&games); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	for i := range games {
		if err := h.games.Insert(context.TODO(), games[i]); err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{})
}

func (h *Internal) saveModule(ctx echo.Context) error {
	var m map[string]interface{}
	if err := ctx.Bind(&m); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	t := store.ModuleType(m["type"].(string))
	var mod store.Module
	switch t {
	case store.ModuleTypeBreaker:
		var b store.Breaker
		mapstructure.Decode(m, &b)
		mod = &b
	case store.ModuleTypeFreeGamesGroup:
		var b store.FreeGamesGroup
		d, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			TagName: "json",
			Result:  &b,
		})
		if err := d.Decode(m); err != nil {
			return err
		}
		mod = &b
	}

	if err := h.storefronts.InsertModule(context.TODO(), mod); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{})
}
