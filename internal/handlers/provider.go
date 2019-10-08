package handlers

import (
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/qilin/crm-api/internal/dispatcher/common"
	validator "gopkg.in/go-playground/validator.v9"
)

func ProviderHandlers(initial config.Initial, validator *validator.Validate, set provider.AwareSet) (common.Handlers, func(), error) {
	hSet := common.HandlerSet{
		Validate: validator,
		AwareSet: set,
	}

	return []common.Handler{
		NewExampleGroup(hSet),
	}, func() {}, nil
}
