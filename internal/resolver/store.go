package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	cacheStore "github.com/eko/gocache/store"
	"github.com/gurukami/typ/v2"
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

func (r *storeQueryResolver) Game(
	ctx context.Context,
	obj *graphql.StoreQuery,
	id string,
) (*store.Game, error) {

	d, e := r.cache.Get("game:" + id)

	if d != nil && e == nil {
		game := &store.Game{}
		e := json.Unmarshal(d.([]byte), game)
		if e != nil {
			return nil, e
		}
		return game, nil
	}

	game, e := r.repo.Games.Get(ctx, id)
	if e != nil {
		return nil, e
	}

	b, e := json.Marshal(game)
	if e != nil {
		return nil, e
	}

	e = r.cache.Set("game:"+id, b, &cacheStore.Options{Cost: 2})
	if e != nil {
		return nil, e
	}

	return game, nil
}

func (r *storeQueryResolver) Games(
	ctx context.Context,
	obj *graphql.StoreQuery,
	id *string,
	genre *store.Genre,
	top *int,
) ([]*store.Game, error) {

	key := fmt.Sprintf("games:%v:%v:%v", typ.Of(id).String().V(), typ.Of(genre).String().V(), typ.Of(top).String().V())

	d, e := r.cache.Get(key)

	if d != nil && e == nil {
		var games []*store.Game
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

	if id != nil {
		games = filter(games, func(g *store.Game) bool {
			return g.ID == *id
		})
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

func filter(games []*store.Game, matcher func(*store.Game) bool) []*store.Game {
	var res = make([]*store.Game, 0, len(games))
	for _, g := range games {
		if matcher(g) {
			res = append(res, g)
		}
	}
	return res
}
