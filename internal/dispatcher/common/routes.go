package common

import (
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type Groups struct {
	Auth    *echo.Group
	GraphQL *echo.Group
	V1      *echo.Group
	Common  *echo.Echo
}

// Handler
type Handler interface {
	Route(groups *Groups)
}

type Handlers []Handler

// Validate
type Validator interface {
	Use(validator *validator.Validate)
}

// HandlerSet
type HandlerSet struct {
	Validate *validator.Validate
	AwareSet provider.AwareSet
}
