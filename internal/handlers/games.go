package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/qilin/crm-api/internal/db/domain/store"
	"github.com/qilin/crm-api/internal/db/repo"
	"github.com/qilin/crm-api/internal/dispatcher/common"
)

type Internal struct {
	games *repo.GamesRepo
}

func NewInternal(games *repo.GamesRepo) *Internal {
	return &Internal{
		games: games,
	}
}

func (h *Internal) Route(groups *common.Groups) {
	groups.V1.POST("/internal/games", h.publishGames)
}

func (h *Internal) publishGames(ctx echo.Context) error {
	var games []store.Game
	if err := ctx.Bind(&games); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	for i := range games {
		if err := h.games.Insert(context.TODO(), &games[i]); err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{})
}
