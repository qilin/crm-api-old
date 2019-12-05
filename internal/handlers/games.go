package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/qilin/crm-api/internal/db/domain/store"
	"github.com/qilin/crm-api/internal/db/repo"
	"github.com/qilin/crm-api/internal/dispatcher/common"
)

type GamesGroup struct {
	games *repo.GamesRepo
}

func NewGamesGroup(games *repo.GamesRepo) *GamesGroup {
	return &GamesGroup{
		games: games,
	}
}

func (h *GamesGroup) Route(groups *common.Groups) {
	groups.V1.POST("/games", func(ctx echo.Context) error {
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
	})

}
