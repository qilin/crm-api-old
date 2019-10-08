package common

import (
	"context"

	"github.com/labstack/echo/v4"
)

const (
	Prefix       = "internal.dispatcher"
	UnmarshalKey = "dispatcher"

	// context keys
	UserContextKey = "user"

	// routes
	AuthGroupPath    = "/auth"
	GraphQLGroupPath = ""
)

func ExtractUserContext(ctx context.Context) *AuthUser {
	if user, ok := ctx.Value(UserContextKey).(*AuthUser); ok {
		return user
	}
	return &AuthUser{}
}

func SetUserContext(ctx echo.Context, user *AuthUser) {
	ctx.Set(UserContextKey, user)
}

// AuthUser
type AuthUser struct {
	Id    int
	Email string
	Roles map[string]bool
}

func (a *AuthUser) IsEmpty() bool {
	return len(a.Email) == 0
}
