package resolver

import (
	"context"

	"github.com/qilin/crm-api/internal/auth"

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

func (r *authMutationResolver) PasswordUpdate(ctx context.Context, obj *graphql1.AuthMutation, old string, new string) (*graphql1.PasswordUpdateResponse, error) {
	panic("not implemented")
}

func (r *authQueryResolver) Profile(ctx context.Context, obj *graphql1.AuthQuery) (*graphql1.User, error) {
	user := auth.ExtractUserContext(ctx)
	u, err := r.repo.Users.Get(ctx, user.Id)
	if err != nil {
		return nil, err
	}
	return &graphql1.User{
		ID:        u.ID,
		Status:    string(u.Status),
		Email:     u.Email,
		Phone:     u.Phone,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		PhotoURL:  u.PhotoURL,
		Language:  u.Language,
	}, nil
}

func (r *mutationResolver) Auth(ctx context.Context) (*graphql1.AuthMutation, error) {
	return &graphql1.AuthMutation{}, nil
}

func (r *queryResolver) Auth(ctx context.Context) (*graphql1.AuthQuery, error) {
	return &graphql1.AuthQuery{}, nil
}
