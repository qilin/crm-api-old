package resolver

import (
	"context"
	"sort"

	"github.com/qilin/crm-api/internal/db/domain/store"
	"github.com/qilin/crm-api/internal/generated/graphql"
)

func (r *queryResolver) Store(ctx context.Context) (*graphql.StoreQuery, error) {
	return &graphql.StoreQuery{}, nil
}

type storeQueryResolver struct{ *Resolver }

func (r *Resolver) StoreQuery() graphql.StoreQueryResolver {
	return &storeQueryResolver{r}
}

func (r *storeQueryResolver) Game(ctx context.Context, obj *graphql.StoreQuery, id string) (*store.Game, error) {
	return r.repo.Games.Get(ctx, id)
}

func (r *storeQueryResolver) Games(
	ctx context.Context,
	obj *graphql.StoreQuery,
	genre *store.Genre,
	top *int,
) ([]*store.Game, error) {
	games, err := r.repo.Games.All(ctx)
	if err != nil {
		return nil, err
	}

	if genre != nil {
		games = filter(games, func(g *store.Game) bool {
			return g.Genre == store.Genre(*genre)
		})
	}

	if top != nil {
		sort.Slice(games, func(i, j int) bool {
			return games[i].Rating > games[j].Rating
		})
		if len(games) > *top {
			games = games[:*top]
		}
	}

	return games, nil
}

func filter(games []*store.Game, matcher func(*store.Game) bool) []*store.Game {
	var res = make([]*store.Game, 0, len(games))
	for _, g := range games {
		if matcher(g) {
			res = append(res, g)
		}
	}
	return res
}

func (r *storeQueryResolver) Module(ctx context.Context, obj *graphql.StoreQuery, id string, locale *string) (store.Module, error) {
	return r.repo.Storefronts.GetModule(ctx, id, store.UserCategoryUnknown)
}

func (r *storeQueryResolver) StoreFront(ctx context.Context, obj *graphql.StoreQuery, locale *string) (*store.StoreFront, error) {
	return &store.StoreFront{}, nil
}
