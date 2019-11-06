package qilin

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/pascaldekloe/jwt"
	"github.com/qilin/crm-api/internal/sdk/common"
)

type QilinAuth struct {
	//
}

func NewQilinAuthenticator() *QilinAuth {
	return &QilinAuth{}
}

func (a QilinAuth) Auth(ctx context.Context, request common.AuthRequest, token *jwt.Claims, log logger.Logger) (response common.AuthResponse, err error) {
	// todo: get token from request.URL,
	// todo: verify token by iss,
	// todo: extract from token {qilinProductUUID, userID}
	// todo: check qilinProductUUID in db or create, extract
	// request.URL

	log.Info("ADAPTER GOT REQUEST")

	return common.AuthResponse{
		Meta: request.Meta,
	}, nil

}

func (a QilinAuth) Order(ctx context.Context, request common.OrderRequest, log logger.Logger) (response common.OrderResponse, err error) {
	return common.OrderResponse{}, nil

}
