package auth

import "context"

type AuthService interface {
	IsAuthenticated(ctx context.Context) error
}

type AuthSrv struct {
	//
}

func (a *AuthSrv) IsAuthenticated(ctx context.Context) error {
	return nil
}

type Authenticator interface {
	SignIn(ctx context.Context) (User, error)
	SignOut(ctx context.Context) error
}
