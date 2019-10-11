package dispatcher

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/qilin/crm-api/internal/dispatcher/common"
	"github.com/qilin/crm-api/internal/generated/graphql"
)

// GetUserDetailsMiddleware
func (d *Dispatcher) graphqlJWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token := ctx.Request().Header.Get(echo.HeaderAuthorization)
		if token == "" {
			d.L().Debug("Can't get Authorization Header")
			return next(ctx)
		}

		claims, err := d.appSet.JwtVerifier.Check(token)
		if err != nil {
			d.L().Debug("JWT Check: " + err.Error())
			return next(ctx)
		}

		email, ok := claims.Set["email"].(string)
		if !ok {
			d.L().Info("can't get email from jwt")
			return next(ctx)
		}

		// create internal session with JWT.id and mapped internal user with ID and roles?

		ctx.Set(common.UserContextKey, &common.AuthUser{
			Email: email,
			// todo: extract roles from somewhere
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
