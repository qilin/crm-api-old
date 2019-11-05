package main

import (
	"context"
	"reflect"

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
		val, ok := ctx.Value("value").(string)
		if !ok {
			val = "default value"
		}

		qilinProductUUID, ok := token.String("qilinProductUUID")
		userID, ok := token.String("userID")

		var meta interface{}
		tof := reflect.TypeOf(request.Meta)
		switch tof.Kind() {
		case reflect.Map:
			meta := request.Meta.(map[string]string)
			meta["qilinProductUUID"] = qilinProductUUID
			meta["userID"] = userID
		case reflect.String:
			meta = val + "?qilinProductUUID=" + qilinProductUUID + "&userID=" + userID
		}
		//if authenticate == nil {
		return common.AuthResponse{
			Token: "plugin token",
			Meta:  meta,
		}, nil
		//}
		//return authenticate(ctx, request, token, log)
	}
}
