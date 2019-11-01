package resolver

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"sort"
	"strconv"

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
	genre *graphql.Genres,
	top *int,
	newest *int,
) ([]*graphql.Game, error) {
	games, err := loadGames()
	if err != nil {
		return nil, err
	}

	if id != nil {
		games = filter(games, func(g *graphql.Game) bool {
			return g.ID == *id
		})
	}

	if genre != nil {
		games = filter(games, func(g *graphql.Game) bool {
			return g.Genre == *genre
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

func loadGames() ([]*graphql.Game, error) {
	data, err := ioutil.ReadFile("./configs/games.json")
	if err != nil {
		return nil, err
	}

	var games []*graphql.Game
	if err := json.Unmarshal(data, &games); err != nil {
		return nil, err
	}
	return games, nil
}

func filter(games []*graphql.Game, matcher func(*graphql.Game) bool) []*graphql.Game {
	var res = make([]*graphql.Game, 0, len(games))
	for _, g := range games {
		if matcher(g) {
			res = append(res, g)
		}
	}
	return res
}
