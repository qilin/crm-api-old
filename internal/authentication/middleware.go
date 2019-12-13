package authentication

import (
	"context"

	"github.com/labstack/echo/v4"
)

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

func IsAuthenticated(ctx context.Context) bool {
	u, ok := ctx.Value(userCtxKey).(User)
	if !ok {
		return false
	}
	if u.IsEmpty() {
		return false
	}
	return true
}
