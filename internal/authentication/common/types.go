package common

import "github.com/labstack/echo/v4"

type AuthenticationProvider interface {
	Provider() string
	SignIn(ctx echo.Context) (user *ExternalUser, url string, err error)
	Callback(ctx echo.Context) (user *ExternalUser, err error)
}
