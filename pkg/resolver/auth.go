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
	// validate email
	// todo: bind validate (wire)
	//e := r.validate.Var(email, "email,required")
	//if e != nil {
	//	return &graphql1.SignupOut{
	//		Status: graphql1.SignupOutStatusBadRequest,
	//	}, e
	//}
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
	//create user
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
	// validate email
	// todo
	user, e := r.repo.User.Get(ctx, email, password)
	if e != nil {
		return &graphql1.SigninOut{
			Status: graphql1.SigninOutStatusServerInternalError,
		}, e
	}
	// gen JWT
	// todo
	jwt := user.Email
	//
	return &graphql1.SigninOut{
		Status: graphql1.SigninOutStatusServerInternalError,
		Token:  jwt,
	}, nil
}

func (r *mutationResolver) Auth(ctx context.Context) (*graphql1.AuthMutation, error) {
	return &graphql1.AuthMutation{}, nil
}

func (r *queryResolver) Auth(ctx context.Context) (*graphql1.AuthQuery, error) {
	return &graphql1.AuthQuery{}, nil
}
