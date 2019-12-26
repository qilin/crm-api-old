package authentication

import (
	"context"

	"github.com/qilin/crm-api/internal/authentication/common"

	"github.com/labstack/echo/v4"
)

func ExtractUserContext(ctx context.Context) *common.User {
	if user, ok := ctx.Value(userCtxKey).(*common.User); ok {
		return user
	}
	return &common.User{}
}

// SetUserContext sets user context into http.Request
// we need it because gqlgen don't knows about echo.Context and uses http.Request context
// but there is no easy way to set Request context.
func SetUserContext(ctx echo.Context, user *common.User) {
	r := ctx.Request()
	newctx := context.WithValue(r.Context(), userCtxKey, user)
	ctx.SetRequest(r.WithContext(newctx))
}

func IsAuthenticated(ctx context.Context) bool {
	u, ok := ctx.Value(userCtxKey).(*common.User)
	if !ok {
		return false
	}
	if u.IsEmpty() {
		return false
	}
	return true
}
