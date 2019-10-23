package auth

import (
	"context"
	"strconv"
	"strings"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// context key
type contextKey struct{ name string }

var userCtxKey = &contextKey{"user"}

func ExtractUserContext(ctx context.Context) *AuthUser {
	if user, ok := ctx.Value(userCtxKey).(*AuthUser); ok {
		return user
	}
	return &AuthUser{}
}

// SetUserContext sets user context into http.Request
// we need it because gqlgen don't knows about echo.Context and uses http.Request context
// but there is no easy way to set Request context.
func SetUserContext(ctx echo.Context, user *AuthUser) {
	r := ctx.Request()
	newctx := context.WithValue(r.Context(), userCtxKey, user)
	ctx.SetRequest(r.WithContext(newctx))
}

// AuthUser
type AuthUser struct {
	Id    int
	Roles map[string]bool
}

func (u *AuthUser) IsEmpty() bool {
	return u.Id == 0
}

// Middleware returns authorization middleware for http server
func (a *Auth) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authHeader := ctx.Request().Header.Get(echo.HeaderAuthorization)
		if authHeader == "" {
			a.log.Error("no auth header")
			return next(ctx)
		}

		rawToken, ok := ExtractTokenFromAuthHeader(authHeader)
		if !ok {
			a.log.Error("invalid auth header")
			return next(ctx)
		}

		var claims = &TokenClaims{}
		// validate jwt token
		if _, err := jwt.ParseWithClaims(rawToken, claims, func(*jwt.Token) (interface{}, error) {
			return a.jwtKeys.Public, nil
		}); err != nil {
			a.log.Error("invalid jwt token: %v", logger.Args(err))
			return next(ctx)
		}

		// create internal session with JWT.id and mapped internal user with ID and roles?
		userId, err := strconv.Atoi(claims.UserID)
		if err != nil {
			a.log.Error("invalid jwt claims: %v", logger.Args(err))
			return next(ctx)
		}

		a.log.Debug("auth user: %d", logger.Args(userId))
		roles := make(map[string]bool)
		for _, r := range claims.AllowedRoles {
			roles[strings.ToUpper(r)] = true
		}

		SetUserContext(ctx, &AuthUser{
			Id:    userId,
			Roles: roles,
		})
		a.log.Debug("auth user: %d", logger.Args(ExtractUserContext(ctx.Request().Context())))

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