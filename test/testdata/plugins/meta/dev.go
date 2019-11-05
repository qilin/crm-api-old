package main

import (
	"context"

	"github.com/pascaldekloe/jwt"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/qilin/crm-api/internal/sdk/common"
)

type plugin struct {
	//
}

var (
	Plugin plugin
)

func (p *plugin) Name() string {
	return "example.plugin"
}

func (p *plugin) Auth(authenticate common.Authenticate) common.Authenticate {
	return func(ctx context.Context, request common.AuthRequest, token *jwt.Claims, log logger.Logger) (response common.AuthResponse, err error) {
		qilinProductUUID, ok := token.String("qilinProductUUID")
		if !ok {
			qilinProductUUID = "unknown"
		}
		userID, ok := token.String("userID")
		if !ok {
			userID = "unknown"
		}

		meta := map[string]string{}
		meta["mode"] = "dev"
		meta["qilinProductUUID"] = qilinProductUUID
		meta["userID"] = userID

		//if authenticate == nil {
		//	return authenticate(ctx, request, token, log)
		//}
		return common.AuthResponse{
			Token: "_jwt_plugin_generated_by_plugin_",
			Meta:  meta,
		}, nil
	}
}
