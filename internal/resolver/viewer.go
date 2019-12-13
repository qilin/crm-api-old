package resolver

import (
	"context"

	"github.com/qilin/crm-api/internal/generated/graphql"
)

func (r *queryResolver) Viewer(ctx context.Context) (*graphql.ViewerQuery, error) {
	return &graphql.ViewerQuery{}, nil
}
