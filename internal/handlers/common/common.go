package common

import (
	"github.com/labstack/echo/v4"
)

// Handler
type Handler interface {
	Route(groups *Groups)
}

type Handlers []Handler

// Handler Routes
type Groups struct {
	Auth    *echo.Group
	GraphQL *echo.Group
	V1      *echo.Group
	Common  *echo.Echo
	SDK     *echo.Group
}
