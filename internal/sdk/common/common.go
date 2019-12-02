package common

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
	"github.com/qilin/crm-api/internal/db/domain"
)

const (
	Prefix                   = "pkg.sdk"
	UnmarshalKey             = "sdk"
	UnmarshalKeyPluginConfig = "sdk.plugin"
)

type SDKMode string

const (
	// todo: rename string const here and in configs
	StoreMode SDKMode = "parent"
	// todo: rename string const here and in configs
	ProviderMode SDKMode = "dev"
)

const (
	QilinProductUUID = "qilinProductUUID"
	UserID           = "userID"
)

// auth
type Authenticate func(ctx context.Context, request AuthRequest, token *jwt.Claims, log logger.Logger) (response AuthResponse, err error)

// order
type Order func(ctx context.Context, request OrderRequest, log logger.Logger) (response OrderResponse, err error)

// SDK
type SDK interface {
	Authenticate(ctx context.Context, request AuthRequest, token *jwt.Claims, log logger.Logger) (response AuthResponse, err error)
	GetProductByUUID(uuid string) (*domain.StoreGamesItem, error)
	IframeHtml(qiliProductUUID string) (string, error)
	IssueJWT(userId, qilinProductUUID string) ([]byte, error)
	Mode() SDKMode
	MapExternalUserToUser(iss string, externalId string) (string, error)
	Order(ctx context.Context, request OrderRequest, log logger.Logger) (response OrderResponse, err error)
	PluginsRoute(echo *echo.Echo)
	Verify(token []byte) (*jwt.Claims, error)
	ActionsLog() domain.ActionsLog
}
