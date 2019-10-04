package resolver

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/qilin/crm-api/internal/dispatcher/common"

	"github.com/dgrijalva/jwt-go"
	"github.com/gurukami/typ/v2"

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
	// validate input
	// 1. validate email
	e := r.validate.Var(email, "email,required")
	if e != nil {
		return &graphql1.SignupOut{
			Status: graphql1.SignupOutStatusBadRequest,
		}, nil
	}
	// 2. validate password length
	e = r.validate.Var(email, "password,required")
	if e != nil {
		return &graphql1.SignupOut{
			Status: graphql1.SignupOutStatusBadRequest,
		}, nil
	}

	// check email already taken
	isExists, e := r.repo.User.IsExistsEmail(ctx, email)
	if e != nil {
		return &graphql1.SignupOut{
			Status: graphql1.SignupOutStatusServerInternalError,
		}, e
	}
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
	e := r.validate.Var(email, "email,required")
	if e != nil {
		return &graphql1.SigninOut{
			Status: graphql1.SigninOutStatusBadRequest,
		}, e
	}

	// get user
	user, e := r.repo.User.GetByEmail(ctx, email)
	if e != nil {
		return &graphql1.SigninOut{
			Status: graphql1.SigninOutStatusServerInternalError,
		}, e
	}
	if e := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); e != nil {
		return &graphql1.SigninOut{
			Status: graphql1.SigninOutStatusBadRequest,
		}, nil
	}

	// gen JWT
	jwt, e := GenerateJWT(r.jwt.Alg, []byte(r.jwt.Private), map[string]interface{}{})
	if e != nil {
		return &graphql1.SigninOut{
			Status: graphql1.SigninOutStatusServerInternalError,
		}, e
	}

	//
	return &graphql1.SigninOut{
		Status: graphql1.SigninOutStatusOk,
		Token:  jwt,
	}, nil
}

func (r *authQueryResolver) Me(ctx context.Context, obj *graphql1.AuthQuery) (*graphql1.User, error) {
	u := common.ExtractUserContext(ctx)
	user, e := r.repo.User.Get(ctx, u.Id)
	if e != nil {
		return nil, e
	}
	return &graphql1.User{
		Email: user.Email,
	}, nil
}

func (r *authQueryResolver) Signout(ctx context.Context, obj *graphql1.AuthQuery) (*graphql1.SignoutOut, error) {
	return &graphql1.SignoutOut{
		Status: graphql1.AuthenticatedRequestStatusOk,
	}, nil
}

func (r *mutationResolver) Auth(ctx context.Context) (*graphql1.AuthMutation, error) {
	return &graphql1.AuthMutation{}, nil
}

func (r *queryResolver) Auth(ctx context.Context) (*graphql1.AuthQuery, error) {
	return &graphql1.AuthQuery{}, nil
}

func GenerateJWT(alg string, privateKey []byte, fields map[string]interface{}) (string, error) {
	type CustomClaims struct {
		UserID int64 `json:"user_id"`
		jwt.StandardClaims
	}

	claims := &CustomClaims{}

	for k, v := range fields {
		switch k {
		case "aud":
			claims.Audience = typ.Of(v).String().V()
		case "exp":
			claims.ExpiresAt = typ.Of(v).Int64().V()
		case "id":
			claims.Id = typ.Of(v).String().V()
		case "iat":
			claims.IssuedAt = time.Now().Unix()
		case "iss":
			claims.Issuer = typ.Of(v).String().V()
		case "nbf":
			claims.NotBefore = typ.Of(v).Int64().V()
		case "sub":
			claims.Subject = typ.Of(v).String().V()
		}
	}

	var (
		sPrivateKey    interface{}
		err            error
		signinigMethod jwt.SigningMethod
	)

	switch alg {
	case "ES256":
		signinigMethod = jwt.SigningMethodES256
		sPrivateKey, err = jwt.ParseECPrivateKeyFromPEM(privateKey)
		if err != nil {
			return "", err
		}
	case "RS256":
		signinigMethod = jwt.SigningMethodRS256
		sPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKey)
		if err != nil {
			return "", err
		}
	}

	token := jwt.NewWithClaims(signinigMethod, claims)
	ss, err := token.SignedString(sPrivateKey)

	if err != nil {
		return "", err
	}

	return ss, nil
}
