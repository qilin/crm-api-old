package handlers

import (
	"net/http"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
	"github.com/qilin/crm-api/internal/dispatcher/common"
)

const (
	examplePath = "/example"
)

type ExampleGroup struct {
	dispatch common.HandlerSet
	provider.LMT
}

func NewExampleGroup(set common.HandlerSet) *ExampleGroup {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "PriceGroup"})
	return &ExampleGroup{
		dispatch: set,
		LMT:      &set.AwareSet,
	}
}

func (h *ExampleGroup) Route(groups *common.Groups) {
	groups.Common.GET(examplePath, h.getExample)
}

// Get currency and region by country code
// GET /api/v1/price_group/country
func (h *ExampleGroup) getExample(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "example")
}
