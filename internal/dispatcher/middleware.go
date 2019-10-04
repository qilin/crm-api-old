package dispatcher

import (
	"regexp"

	"github.com/qilin/crm-api/internal/dispatcher/common"

	"github.com/labstack/echo/v4"
)

var tokenRegex = regexp.MustCompile("Bearer ([A-z0-9_.-]{10,})")

// GetUserDetailsMiddleware
func (d *Dispatcher) GetUserDetailsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		auth := ctx.Request().Header.Get(echo.HeaderAuthorization)
		if auth == "" {
			return next(ctx)
		}

		match := tokenRegex.FindStringSubmatch(auth)
		if len(match) < 1 {
			return next(ctx)
		}

		// @todo parse jwt and fill
		ctx.Set("user", common.AuthUser{
			Id: 1,
		})

		return next(ctx)
	}
}
