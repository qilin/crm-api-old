package resolver

import (
	"context"
	"encoding/json"
	"fmt"

	cacheStore "github.com/eko/gocache/store"
	"github.com/gurukami/typ/v2"
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

func (r *storeQueryResolver) Game(ctx context.Context, obj *graphql.StoreQuery, id string) (store.Game, error) {
	game, err := r.repo.Games.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (r *storeQueryResolver) GameBySlug(ctx context.Context, obj *graphql.StoreQuery, slug string) (store.Game, error) {
	game, err := r.repo.Games.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (r *storeQueryResolver) Games(
	ctx context.Context,
	obj *graphql.StoreQuery,
	genre *store.Genre,
) ([]store.Game, error) {

	key := fmt.Sprintf("games:%v", typ.Of(genre).String().V())
	d, e := r.cache.Get(key)
	if d != nil && e == nil {
		var games store.GameSlice
		e := json.Unmarshal(d.([]byte), &games)
		if e != nil {
			return nil, e
		}
		return games, nil
	}

	games, err := r.repo.Games.All(ctx)
	if err != nil {
		return nil, err
	}

	if genre != nil {
		games = filter(games, func(g store.Game) bool {
			for _, g := range g.Common().Genres {
				if g == *genre {
					return true
				}
			}
			return false
		})
	}

	b, e := json.Marshal(games)
	if e != nil {
		return nil, e
	}

	e = r.cache.Set(key, b, &cacheStore.Options{Cost: 2})
	if e != nil {
		return nil, e
	}

	return games, nil
}

func filter(games []store.Game, matcher func(store.Game) bool) []store.Game {
	var res = make([]store.Game, 0, len(games))
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

type freeGameOfferResolver struct{ *Resolver }

func (r *Resolver) FreeGameOffer() graphql.FreeGameOfferResolver {
	return &freeGameOfferResolver{r}
}

func (r *freeGameOfferResolver) Game(ctx context.Context, obj *store.FreeGameOffer) (store.Game, error) {
	return r.repo.Games.Get(ctx, obj.GameID)
}
