package common

import "context"

const (
	Prefix                 = "internal.dispatcher"
	UnmarshalKey           = "dispatcher"
	UnmarshalAuthConfigKey = "dispatcher.auth"
)

func ExtractUserContext(ctx context.Context) *AuthUser {
	if user, ok := ctx.Value("user").(*AuthUser); ok {
		return user
	}
	return &AuthUser{}
}

// AuthUser
type AuthUser struct {
	Id    int
	Email string
	Roles map[string]bool
}

func (a *AuthUser) IsEmpty() bool {
	return a.Id == 0
}
