package resolver

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/qilin/crm-api/internal/generated/graphql"
)

func (r *queryResolver) Store(ctx context.Context) (*graphql.StoreQuery, error) {
	data, err := ioutil.ReadFile("./configs/games.json")
	if err != nil {
		return nil, err
	}

	var games []*graphql.Game
	if err := json.Unmarshal(data, &games); err != nil {
		return nil, err
	}

	return &graphql.StoreQuery{
		Games: games,
	}, nil
}
