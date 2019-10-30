package sdk

import (
	"net/http"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
	"github.com/qilin/crm-api/internal/dispatcher/common"
)

const (
	sdkAuthRoute   = "/auth"
	sdkOrderRoute  = "/order"
	sdkHealthRoute = "/health"
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
	h.L().Info("create routes")
	groups.SDK.POST(sdkAuthRoute, h.postAuth)
	groups.SDK.POST(sdkOrderRoute, h.postOrder)
	groups.SDK.POST(sdkHealthRoute, h.getHealth)
}

// POST /sdk/v1/auth
func (h *SDKGroup) postAuth(ctx echo.Context) error {
	r := new(AuthRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code: errBadRequest,
			Msg:  StatusText(errBadRequest),
		})
	}

	if len(r.URL) == 0 {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code: errAuthRequestURLEmpty,
			Msg:  StatusText(errAuthRequestURLEmpty),
		})
	}

	return ctx.JSON(http.StatusOK, AuthResponse{
		Meta: r.Meta.(string),
	})
}

// POST /sdk/v1/order
func (h *SDKGroup) postOrder(ctx echo.Context) error {
	r := new(OrderRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code: errBadRequest,
			Msg:  StatusText(errBadRequest),
		})
	}

	if len(r.Data) == 0 {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code: errOrderRequestDataEmpty,
			Msg:  StatusText(errOrderRequestDataEmpty),
		})
	}

	return ctx.JSON(http.StatusOK, OrderResponse{
		//Data: r.Data,
	})
}

func (h *SDKGroup) getHealth(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}
