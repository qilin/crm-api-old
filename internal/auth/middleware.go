package auth

import (
	"context"
	"strconv"
	"strings"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/labstack/echo/v4"
)

// context key
type contextKey struct{ name string }

var userCtxKey = &contextKey{"user"}

func ExtractUserContext(ctx context.Context) *User {
	if user, ok := ctx.Value(userCtxKey).(*User); ok {
		return user
	}
	return &User{}
}

// SetUserContext sets user context into http.Request
// we need it because gqlgen don't knows about echo.Context and uses http.Request context
// but there is no easy way to set Request context.
func SetUserContext(ctx echo.Context, user *User) {
	r := ctx.Request()
	newctx := context.WithValue(r.Context(), userCtxKey, user)
	ctx.SetRequest(r.WithContext(newctx))
}

// User
type User struct {
	Id    int
	Roles map[string]bool
}

func (u *User) IsEmpty() bool {
	return u.Id == 0
}

// Middleware returns authorization middleware for http server
func (a *Auth) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		rawToken, ok := a.extractToken(ctx)
		if !ok {
			a.L().Error("invalid auth header")
			return next(ctx)
		}

		// validate jwt token
		var claims = &AccessTokenClaims{}
		if err := a.jwtKeys.Parse(rawToken, claims); err != nil {
			a.L().Error("invalid jwt token: %v", logger.Args(err))
			return next(ctx)
		}

		// create internal session with JWT.id and mapped internal user with ID and roles?
		userId, err := strconv.Atoi(claims.UserID)
		if err != nil {
			a.L().Error("invalid jwt claims: %v", logger.Args(err))
			return next(ctx)
		}

		a.L().Debug("auth user: %d", logger.Args(userId))
		roles := make(map[string]bool)
		roles[strings.ToUpper(claims.Role)] = true

		SetUserContext(ctx, &User{
			Id:    userId,
			Roles: roles,
		})
		a.L().Debug("auth user: %d", logger.Args(ExtractUserContext(ctx.Request().Context())))

		return next(ctx)
	}
}

const bearer = "bearer"

func ExtractTokenFromAuthHeader(ctx echo.Context) (token string, ok bool) {
	authHeader := ctx.Request().Header.Get(echo.HeaderAuthorization)
	if authHeader == "" {
		return "", false
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || !strings.EqualFold(strings.ToLower(authHeaderParts[0]), bearer) {
		return "", false
	}

	return authHeaderParts[1], true
}

func (a *Auth) extractTokenFromCookie(ctx echo.Context) (token string, ok bool) {
	cookie, err := ctx.Request().Cookie(a.cfg.SessionCookieName)
	if err != nil {
		return "", false
	}
	token = cookie.Value
	if token != "" {
		ok = true
	}
	return token, ok
}

func (a *Auth) extractToken(ctx echo.Context) (token string, ok bool) {
	token, ok = ExtractTokenFromAuthHeader(ctx)
	if token != "" {
		return token, ok
	}
	return a.extractTokenFromCookie(ctx)
}
