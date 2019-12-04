package resolver

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/qilin/crm-api/internal/db/domain"

	graphql1 "github.com/qilin/crm-api/internal/generated/graphql"
)

type authMutationResolver struct{ *Resolver }
type authQueryResolver struct{ *Resolver }

func (r *Resolver) AuthMutation() graphql1.AuthMutationResolver {
	return &authMutationResolver{r}
}
func (r *Resolver) AuthQuery() graphql1.AuthQueryResolver {
	return &authQueryResolver{r}
}

func (r *authMutationResolver) SignUp(ctx context.Context, obj *graphql1.AuthMutation, email string, password string) (*graphql1.SignUpResponse, error) {
	// validate input
	// 1. validate email
	e := r.validate.Var(email, "email,required")
	if e != nil {
		return &graphql1.SignUpResponse{
			Status: graphql1.SignUpResponseStatusBadRequest,
		}, nil
	}
	// 2. validate password length
	e = r.validate.Var(email, "password,required")
	if e != nil {
		return &graphql1.SignUpResponse{
			Status: graphql1.SignUpResponseStatusBadRequest,
		}, nil
	}

	// check email already taken
	isExists, e := r.repo.User.IsExistsEmail(ctx, email)
	if e != nil {
		return &graphql1.SignUpResponse{
			Status: graphql1.SignUpResponseStatusServerInternalError,
		}, e
	}
	if isExists {
		return &graphql1.SignUpResponse{
			Status: graphql1.SignUpResponseStatusUserExists,
		}, nil
	}

	//create user
	user := &domain.UserItem{
		Email:    email,
		Password: password,
	}
	e = r.repo.User.Create(ctx, user)
	if e != nil {
		return &graphql1.SignUpResponse{
			Status: graphql1.SignUpResponseStatusServerInternalError,
		}, e
	}
	return &graphql1.SignUpResponse{
		Status: graphql1.SignUpResponseStatusOk,
	}, nil
}

func (r *authMutationResolver) PasswordUpdate(ctx context.Context, obj *graphql1.AuthMutation, old string, new string) (*graphql1.PasswordUpdateResponse, error) {
	panic("not implemented")
}

func (r *authQueryResolver) SignIn(ctx context.Context, obj *graphql1.AuthQuery, email string, password string) (*graphql1.SignInResponse, error) {
	// validate email
	e := r.validate.Var(email, "email,required")
	if e != nil {
		return &graphql1.SignInResponse{
			Status: graphql1.RequestStatusBadRequest,
		}, e
	}

	// get user
	user, e := r.repo.User.FindByEmail(ctx, email)
	if e != nil {
		return &graphql1.SignInResponse{
			Status: graphql1.RequestStatusServerInternalError,
		}, e
	}
	if e := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); e != nil {
		return &graphql1.SignInResponse{
			Status: graphql1.RequestStatusBadRequest,
		}, nil
	}

	jwt := "jwt.token"

	//
	return &graphql1.SignInResponse{
		Status: graphql1.RequestStatusOk,
		Token:  jwt,
	}, nil
}

func (r *authQueryResolver) SignOut(ctx context.Context, obj *graphql1.AuthQuery) (*graphql1.SignOutResponse, error) {
	panic("not implemented")
}

func (r *authQueryResolver) Profile(ctx context.Context, obj *graphql1.AuthQuery) (*graphql1.User, error) {
	panic("not implemented")
}

func (r *mutationResolver) Auth(ctx context.Context) (*graphql1.AuthMutation, error) {
	return &graphql1.AuthMutation{}, nil
}

func (r *queryResolver) Auth(ctx context.Context) (*graphql1.AuthQuery, error) {
	return &graphql1.AuthQuery{}, nil
}
