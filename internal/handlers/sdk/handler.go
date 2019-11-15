package sdk

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
	"github.com/qilin/crm-api/internal/dispatcher/common"
	common2 "github.com/qilin/crm-api/internal/sdk/common"
)

const (
	sdkAuthRoute   = "/auth"
	sdkOrderRoute  = "/order"
	sdkQilinIframe = "/iframe"
	sdkHealthRoute = "/health"
)

type SDKGroup struct {
	sdk common2.SDK
	set common.HandlerSet
	provider.LMT
}

func NewSDKGroup(set common.HandlerSet, sdk common2.SDK) *SDKGroup {
	set.AwareSet.Logger = set.AwareSet.Logger.WithFields(logger.Fields{"router": "SDKGroup"})
	return &SDKGroup{
		sdk: sdk,
		set: set,
		LMT: &set.AwareSet,
	}
}

func (h *SDKGroup) Route(groups *common.Groups) {
	// plugins routes
	h.sdk.PluginsRoute(groups.Common)
	// sdk routes
	groups.SDK.POST(sdkAuthRoute, h.postAuth)
	groups.SDK.POST(sdkOrderRoute, h.postOrder)
	groups.SDK.POST(sdkHealthRoute, h.getHealth)
	groups.SDK.GET(sdkQilinIframe, h.qilinIframe)
}

// POST /sdk/v1/auth
func (h *SDKGroup) postAuth(ctx echo.Context) error {
	r := common2.AuthRequest{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, common2.ErrorResponse{
			Code: errBadRequest,
			Msg:  StatusText(errBadRequest),
		})
	}

	if r.Meta == nil {
		r.Meta = map[string]string{}
	}

	if err := h.set.Validate.Struct(r); err != nil {
		h.L().Info(err.Error())
		return ctx.JSON(http.StatusBadRequest, common2.ErrorResponse{
			Code: errAuthRequestURLEmpty,
			Msg:  StatusText(errAuthRequestURLEmpty),
		})
	}

	var resp common2.AuthResponse
	var err error

	switch h.sdk.Mode() {
	case common2.StoreMode:
		// pass it into plugin
		resp, err = h.parentMode(context.WithValue(context.Background(), "request", ctx.Request()), r)
	case common2.ProviderMode:
		// parse, verify and pass into plugin
		resp, err = h.devMode(context.Background(), r)
	default:
		// parse, verify and pass into adapter
		resp, err = h.qilinMode(context.Background(), r)
	}

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, common2.ErrorResponse{
			Code: errInternalServerError,
			Msg:  err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, resp)
}

// POST /sdk/v1/order
func (h *SDKGroup) postOrder(ctx echo.Context) error {
	r := new(common2.OrderRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, common2.ErrorResponse{
			Code: errBadRequest,
			Msg:  StatusText(errBadRequest),
		})
	}

	if len(r.Data) == 0 {
		return ctx.JSON(http.StatusBadRequest, common2.ErrorResponse{
			Code: errOrderRequestDataEmpty,
			Msg:  StatusText(errOrderRequestDataEmpty),
		})
	}

	// @todo

	return ctx.JSON(http.StatusOK, common2.OrderResponse{
		//Data: r.Data,
	})
}

func (h *SDKGroup) getHealth(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func (h *SDKGroup) qilinIframe(ctx echo.Context) error {
	// todo: extract product uuid from request
	html, err := h.sdk.IframeHtml("")
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, err.Error())
	}
	return ctx.HTML(http.StatusOK, html)
}

func (h *SDKGroup) parentMode(ctx context.Context, r common2.AuthRequest) (common2.AuthResponse, error) {
	claims := &jwt.Claims{}
	return h.sdk.Authenticate(ctx, r, claims, h.L())
}

func (h *SDKGroup) devMode(ctx context.Context, r common2.AuthRequest) (common2.AuthResponse, error) {
	u, err := url.Parse(r.URL)
	if err != nil {
		return common2.AuthResponse{}, err
	}

	claims, err := h.sdk.Verify([]byte(u.Query().Get("jwt")))
	if err != nil {
		return common2.AuthResponse{}, err
	}

	userId, ok := claims.String(common2.UserID)
	if !ok {
		return common2.AuthResponse{}, errors.New(StatusText(errTokenNoUserId))
	}
	qilinProductUUID, ok := claims.String(common2.QilinProductUUID)
	if !ok {
		return common2.AuthResponse{}, errors.New(StatusText(errTokenNoQilinPorductUUID))
	}

	err = h.set.Validate.Var(qilinProductUUID, "uuid4")
	if err != nil {
		return common2.AuthResponse{}, err
	}
	err = h.set.Validate.Var(userId, "min=1")
	if err != nil {
		return common2.AuthResponse{}, err
	}

	// todo: ?

	ctx = context.WithValue(ctx, common2.UserID, userId)
	r.QilinProductUUID = qilinProductUUID

	return h.sdk.Authenticate(ctx, r, claims, h.L())
}

func (h *SDKGroup) qilinMode(ctx context.Context, r common2.AuthRequest) (common2.AuthResponse, error) {
	u, err := url.Parse(r.URL)
	if err != nil {
		return common2.AuthResponse{}, err
	}

	claims, err := h.sdk.Verify([]byte(u.Query().Get("jwt")))
	if err != nil {
		return common2.AuthResponse{}, err
	}

	userId, ok := claims.String(common2.UserID)
	if !ok {
		return common2.AuthResponse{}, errors.New(StatusText(errTokenNoUserId))
	}
	qilinProductUUID, ok := claims.String(common2.QilinProductUUID)
	if !ok {
		return common2.AuthResponse{}, errors.New(StatusText(errTokenNoQilinPorductUUID))
	}

	err = h.set.Validate.Var(qilinProductUUID, "uuid4")
	if err != nil {
		return common2.AuthResponse{}, err
	}
	err = h.set.Validate.Var(userId, "min=1")
	if err != nil {
		return common2.AuthResponse{}, err
	}

	// todo: tmp fix, hardcoded product value; prod needs logic fix inside GetProductByUUID
	// check qilinProductUUID and extract channeling URL
	product, err := h.sdk.GetProductByUUID(qilinProductUUID)
	if err != nil {
		return common2.AuthResponse{}, err
	}

	// todo: check mapping (iss+userID) & userID Qilin
	// todo: tmp fix, hardcoded userID as uuidV4; prod needs logic fix inside MapExternalUserToUser
	userId = h.sdk.MapExternalUserToUser(0, userId)

	// build JWT (qilinProductUUID,userID)
	jwt, err := h.sdk.IssueJWT(userId, product.ID)
	if err != nil {
		return common2.AuthResponse{}, err
	}

	iframe, err := url.Parse(product.URL)
	if err != nil {
		return common2.AuthResponse{}, err
	}

	ifq := iframe.Query()
	ifq.Set("jwt", string(jwt))
	iframe.RawQuery = ifq.Encode()

	// return meta[url] = channeling URL, channeling URL = URL?jwt=<JWT>
	r.Meta = map[string]string{
		"url": iframe.String(),
	}

	return common2.AuthResponse{
		Meta: r.Meta,
	}, nil
}
