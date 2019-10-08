package dispatcher

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/qilin/crm-api/generated/graphql"
	"github.com/qilin/crm-api/internal/dispatcher/common"
)

// GetUserDetailsMiddleware
func (d *Dispatcher) graphqlJWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authHeader := ctx.Request().Header.Get(echo.HeaderAuthorization)
		if authHeader == "" {
			return next(ctx)
		}

		token, ok := ExtractTokenFromAuthHeader(authHeader)
		if !ok {
			return next(ctx)
		}

		claims, err := d.appSet.JwtVerifier.Check(token)
		if err != nil {
			return next(ctx)
		}

		email, ok := claims.Set["email"].(string)
		if !ok {
			return next(ctx)
		}

		// create internal session with JWT.id and mapped internal user with ID and roles?

		ctx.Set(common.UserContextKey, common.AuthUser{
			Email: email,
			Roles: map[string]bool{
				graphql.RoleEnumUser.String(): true,
			},
		})

		return next(ctx)
	}
}

const bearer = "bearer"

func ExtractTokenFromAuthHeader(val string) (token string, ok bool) {
	authHeaderParts := strings.Split(val, " ")
	if len(authHeaderParts) != 2 || !strings.EqualFold(strings.ToLower(authHeaderParts[0]), bearer) {
		return "", false
	}

	return authHeaderParts[1], true
}
