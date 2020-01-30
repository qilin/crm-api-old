package handlers

import (
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/qilin/crm-api/internal/auth"
	"github.com/qilin/crm-api/internal/handlers/common"
	"github.com/qilin/crm-api/internal/webhooks"
	"github.com/qilin/crm-api/pkg/graphql"
	"gopkg.in/go-playground/validator.v9"
)

type HandlerSet struct {
	GraphQL  *graphql.GraphQL
	WebHooks *webhooks.WebHooks
	Auth     *auth.Auth
	Internal *Internal
}

// ProviderHandlers
func ProviderHandlers(initial config.Initial, validator *validator.Validate, set provider.AwareSet, handlers HandlerSet) (common.Handlers, func(), error) {
	listHandlers := []common.Handler{
		handlers.GraphQL,
		handlers.WebHooks,
		handlers.Auth,
		handlers.Internal,
	}
	return listHandlers, func() {}, nil
}
