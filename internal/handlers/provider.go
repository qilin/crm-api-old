package handlers

import (
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/qilin/crm-api/internal/auth"
	"github.com/qilin/crm-api/internal/dispatcher/common"
	"github.com/qilin/crm-api/internal/webhooks"
	"github.com/qilin/crm-api/pkg/graphql"
	"gopkg.in/go-playground/validator.v9"
)

type Handlers struct {
	GraphQL  *graphql.GraphQL
	WebHooks *webhooks.WebHooks
	Auth     *auth.Auth
}

// ProviderHandlers
func ProviderHandlers(initial config.Initial, validator *validator.Validate, set provider.AwareSet, handlers Handlers) (common.Handlers, func(), error) {
	hSet := common.HandlerSet{
		Validate: validator,
		AwareSet: set,
	}
	listHandlers := []common.Handler{
		NewExampleGroup(hSet),
		handlers.GraphQL,
		handlers.WebHooks,
		handlers.Auth,
	}
	return listHandlers, func() {}, nil
}
