package resolver

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"sort"
	"strconv"

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

func (r *storeQueryResolver) Games(
	ctx context.Context,
	obj *graphql.StoreQuery,
	id *string,
	genre *store.Genre,
	top *int,
	newest *int,
) ([]*store.Game, error) {
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

	if newest != nil {
		sort.Slice(games, func(i, j int) bool {
			iid, _ := strconv.Atoi(games[i].ID)
			jid, _ := strconv.Atoi(games[j].ID)
			return iid > jid
		})
		if len(games) > *newest {
			games = games[:*newest]
		}
	}

	return games, nil
}

func loadGames() ([]*store.Game, error) {
	data, err := ioutil.ReadFile("./configs/games.json")
	if err != nil {
		return nil, err
	}

	var games []*store.Game
	if err := json.Unmarshal(data, &games); err != nil {
		return nil, err
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
