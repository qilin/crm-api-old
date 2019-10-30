package common

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/pascaldekloe/jwt"
)

const (
	Prefix       = "pkg.sdk"
	UnmarshalKey = "sdk"
)

type SDKMode string

const (
	ParentMode    SDKMode = "parent"
	DeveloperMode SDKMode = "dev"
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
	Mode() SDKMode
	Verify(token []byte) (*jwt.Claims, error)
	Authenticate(ctx context.Context, request AuthRequest, token *jwt.Claims, log logger.Logger) (response AuthResponse, err error)
	Order(ctx context.Context, request OrderRequest, log logger.Logger) (response OrderResponse, err error)
}
