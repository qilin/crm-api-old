package handlers

import (
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/qilin/crm-api/internal/dispatcher/common"
	"github.com/qilin/crm-api/internal/handlers/sdk"
	"github.com/qilin/crm-api/internal/webhooks"
	"github.com/qilin/crm-api/pkg/graphql"
	"gopkg.in/go-playground/validator.v9"
)

// ProviderHandlers
func ProviderHandlers(initial config.Initial, validator *validator.Validate, set provider.AwareSet, ql *graphql.GraphQL, wh *webhooks.WebHooks) (common.Handlers, func(), error) {
	return []common.Handler{
		ql,
		wh,
	}, func() {}, nil
}

// ProviderHandlers
func ProviderSDKHandlers(validator *validator.Validate, set provider.AwareSet) (common.Handlers, func(), error) {
	hSet := common.HandlerSet{
		Validate: validator,
		AwareSet: set,
	}
	return []common.Handler{
		sdk.NewSDKGroup(hSet),
	}, func() {}, nil
}
