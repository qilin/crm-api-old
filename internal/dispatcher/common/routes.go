package common

import (
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"gopkg.in/go-playground/validator.v9"
)

// Validate
type Validator interface {
	Use(validator *validator.Validate)
}

// HandlerSet
type HandlerSet struct {
	Validate *validator.Validate
	AwareSet provider.AwareSet
}
