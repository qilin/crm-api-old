package sdk

import (
	"net/http"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
	"github.com/qilin/crm-api/internal/dispatcher/common"
)

const (
	sdkAuth = "/auth"
)

type SDKGroup struct {
	dispatch common.HandlerSet
	provider.LMT
}

func NewSDKGroup(set common.HandlerSet) *SDKGroup {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "PriceGroup"})
	return &SDKGroup{
		dispatch: set,
		LMT:      &set.AwareSet,
	}
}

func (h *SDKGroup) Route(groups *common.Groups) {
	groups.SDK.POST(sdkAuth, h.postAuth)
}

// Get currency and region by country code
// POST /sdk/v1/auth
func (h *SDKGroup) postAuth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "",
		Meta: "{'user_id':10}",
	})
}
