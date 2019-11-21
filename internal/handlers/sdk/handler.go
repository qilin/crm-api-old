package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"

	"github.com/qilin/crm-api/internal/dispatcher/common"
	common2 "github.com/qilin/crm-api/internal/sdk/common"
	"github.com/qilin/crm-api/pkg/qilin"
)

const (
	sdkAuthRoute   = "/auth"
	sdkOrderRoute  = "/order"
	sdkItemsRoute  = "/items"
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
	groups.SDK.GET(sdkQilinIframe, h.hubIframe)
	groups.SDK.GET(sdkItemsRoute, h.getItems)
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
	c := context.WithValue(context.Background(), "request", ctx.Request())

	switch h.sdk.Mode() {
	case common2.StoreMode:
		// pass it into plugin
		resp, err = h.storeMode(c, r)
	case common2.ProviderMode:
		// parse, verify and pass into plugin
		resp, err = h.providerMode(c, r)
	default:
		// parse, verify and pass into adapter
		resp, err = h.hubMode(c, r)
	}

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, common2.ErrorResponse{
			Code: errInternalServerError,
			Msg:  err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, resp)
}

// GET /sdk/v1/items
func (h *SDKGroup) getItems(ctx echo.Context) error {
	var gameId = ctx.QueryParam("game_id")
	var itemId = ctx.QueryParam("item_id")

	product, err := h.sdk.GetProductByUUID(gameId)
	if err != nil {
		return err
	}

	u, err := url.Parse(qilin.ItemsURL(product.URL))
	if err != nil {
		return err
	}

	u.RawQuery = url.Values{
		"game_id": []string{gameId},
		"item_id": []string{itemId},
	}.Encode()

	resp, err := http.Get(u.String())
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

	return ctx.JSONBlob(http.StatusOK, d)

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

	product, err := h.sdk.GetProductByUUID(r.GameID)
	if err != nil {
		return err
	}

	// todo: check mapping (iss+userID) & userID Qilin
	// todo: tmp fix, hardcoded userID as uuidV4; prod needs logic fix inside MapExternalUserToUser
	userId := h.sdk.MapExternalUserToUser(0, r.UserID)

	qilin.OrderURL(product.URL)

	req := common2.OrderRequest{
		GameID: r.GameID,
		UserID: userId,
		ItemID: r.ItemID,
	}
	data, err := json.Marshal(&req)
	if err != nil {
		h.L().Error(err.Error())
		return ctx.JSON(http.StatusInternalServerError, common2.ErrorResponse{
			Code: errInternalServerError,
			Msg:  StatusText(errInternalServerError),
		})
	}

	resp, err := http.Post(qilin.OrderURL(product.URL), "application/json;charset=utf-8", bytes.NewReader(data))
	if err != nil {
		h.L().Error(err.Error())
		return ctx.JSON(http.StatusInternalServerError, common2.ErrorResponse{
			Code: errInternalServerError,
			Msg:  StatusText(errInternalServerError),
		})
	}

	if resp.StatusCode != http.StatusOK {
		h.L().Error("request failed")
		return ctx.JSON(http.StatusInternalServerError, common2.ErrorResponse{
			Code: errInternalServerError,
			Msg:  StatusText(errInternalServerError),
		})
	}

	return ctx.JSON(http.StatusOK, common2.OrderResponse{
		//Data: r.Data,
	})
}

func (h *SDKGroup) getHealth(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func (h *SDKGroup) hubIframe(ctx echo.Context) error {
	// todo: extract product uuid from request
	html, err := h.sdk.IframeHtml("")
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, err.Error())
	}
	return ctx.HTML(http.StatusOK, html)
}

func (h *SDKGroup) storeMode(ctx context.Context, r common2.AuthRequest) (common2.AuthResponse, error) {
	claims := &jwt.Claims{}
	return h.sdk.Authenticate(ctx, r, claims, h.L())
}

func (h *SDKGroup) providerMode(ctx context.Context, r common2.AuthRequest) (common2.AuthResponse, error) {
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

func (h *SDKGroup) hubMode(ctx context.Context, r common2.AuthRequest) (common2.AuthResponse, error) {
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

	iframe, err := url.Parse(qilin.IframeURL(product.URL))
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
