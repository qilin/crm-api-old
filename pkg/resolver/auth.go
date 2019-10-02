package resolver

import (
	"context"

	"github.com/qilin/crm-api/internal/db/domain"

	graphql1 "github.com/qilin/crm-api/generated/graphql"
)

type authMutationResolver struct{ *Resolver }
type authQueryResolver struct{ *Resolver }

func (r *Resolver) AuthMutation() graphql1.AuthMutationResolver {
	return &authMutationResolver{r}
}
func (r *Resolver) AuthQuery() graphql1.AuthQueryResolver {
	return &authQueryResolver{r}
}

func (r *authMutationResolver) Signup(ctx context.Context, obj *graphql1.AuthMutation, email string, password string) (*graphql1.SignupOut, error) {
	// validate email and password
	// todo
	// check email already taken
	isExists, e := r.repo.User.IsExistsEmail(ctx, email)
	if e != nil {
		return &graphql1.SignupOut{
			Status: graphql1.SignupOutStatusServerInternalError,
		}, e
	}
	// if email was already taken
	if isExists {
		return &graphql1.SignupOut{
			Status: graphql1.SignupOutStatusUserExists,
		}, nil
	}
	// create user
	user := &domain.UserItem{
		Email:    email,
		Password: password,
	}
	e = r.repo.User.Create(ctx, user)
	if e != nil {
		return &graphql1.SignupOut{
			Status: graphql1.SignupOutStatusServerInternalError,
		}, e
	}
	return &graphql1.SignupOut{
		Status: graphql1.SignupOutStatusOk,
	}, nil
}

func (r *authQueryResolver) Signin(ctx context.Context, obj *graphql1.AuthQuery, email string, password string) (*graphql1.SigninOut, error) {
	return &graphql1.SigninOut{
		Status: graphql1.SigninOutStatusServerInternalError,
	}, nil
}
func (r *authQueryResolver) Me(ctx context.Context, obj *graphql1.AuthQuery) (*graphql1.MeOut, error) {
	return &graphql1.MeOut{
		Status: graphql1.AuthenticatedRequestStatusForbidden,
	}, nil
}

func (r *mutationResolver) Auth(ctx context.Context) (*graphql1.AuthMutation, error) {
	return &graphql1.AuthMutation{}, nil
}

func (r *queryResolver) Auth(ctx context.Context) (*graphql1.AuthQuery, error) {
	return &graphql1.AuthQuery{}, nil
}
