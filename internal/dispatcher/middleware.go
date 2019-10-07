package dispatcher

import (
	"regexp"
	"strings"

	"github.com/qilin/crm-api/generated/graphql"

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

		_, ok := extractTokenFromAuthHeader(auth)
		if !ok {
			return next(ctx)
		}

		// @todo parse jwt and fill
		ctx.Set("user", common.AuthUser{
			Id: 1,
			Roles: map[string]bool{
				graphql.RoleEnumUser.String(): true,
			},
		})

		return next(ctx)
	}
}

const bearer = "bearer"

func extractTokenFromAuthHeader(val string) (token string, ok bool) {
	authHeaderParts := strings.Split(val, " ")
	if len(authHeaderParts) != 2 || !strings.EqualFold(strings.ToLower(authHeaderParts[0]), bearer) {
		return "", false
	}

	return authHeaderParts[1], true
}
