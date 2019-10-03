package resolver

import (
	"context"
	"fmt"

	"gopkg.in/go-playground/validator.v9"

	"github.com/99designs/gqlgen/graphql"
	graphql1 "github.com/qilin/crm-api/generated/graphql"
	"github.com/qilin/crm-api/internal/db/domain"
	"github.com/qilin/crm-api/internal/db/trx"
	gqErrs "github.com/qilin/crm-api/pkg/graphql/errors"
	"github.com/qilin/go-core/invoker"
	"github.com/qilin/go-core/logger"
	"github.com/qilin/go-core/provider"
)

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// Config custom graphql settings resolvers
type Config struct {
	Debug   bool `fallback:"shared.debug"`
	invoker invoker.Invoker
}

// OnReload
func (c *Config) OnReload(callback func(ctx context.Context)) {
	c.invoker.OnReload(callback)
}

// Reload
func (c *Config) Reload(ctx context.Context) {
	c.invoker.Reload(ctx)
}

// Resolver config graphql resolvers
type Resolver struct {
	ctx      context.Context
	cfg      *Config
	repo     Repo
	validate validator.Validate
	trx      *trx.Manager
	provider.LMT
}

// Mutation returns root graphql mutation resolver
func (r *Resolver) Mutation() graphql1.MutationResolver {
	return &mutationResolver{r}
}

// Query returns root graphql query resolver
func (r *Resolver) Query() graphql1.QueryResolver {
	return &queryResolver{r}
}

// AddErrorf is a convenience method for adding an error to the current response
func (r *Resolver) AddDebugErrorf(ctx context.Context, format string, args ...interface{}) {
	if r.cfg.Debug {
		graphql.AddError(ctx, gqErrs.WrapClientErr(fmt.Errorf(format, args...)))
	}
}

// Repo
type Repo struct {
	User domain.UserRepo
	List domain.ListRepo
}

// New returns instance of config graphql resolvers
func New(ctx context.Context, set provider.AwareSet, appSet AppSet, cfg *Config) graphql1.Config {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": Prefix})
	c := graphql1.Config{
		Resolvers: &Resolver{
			ctx:  ctx,
			cfg:  cfg,
			repo: appSet.Repo,
			//validate: appSet.Validate,
			trx: appSet.Trx,
			LMT: &set,
		},
	}
	return c
}
