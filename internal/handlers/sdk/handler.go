package sdk

import (
	"net/http"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
	"github.com/qilin/crm-api/internal/dispatcher/common"
)

const (
	sdkAuth  = "/auth"
	sdkOrder = "/order"
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
	groups.SDK.POST(sdkOrder, h.postOrder)
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
		Meta: r.Meta,
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
		Data: r.Data,
	})
}
