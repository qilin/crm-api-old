package resolver

import (
	"context"

	"github.com/gurukami/typ/v2"
	"github.com/qilin/crm-api/generated/api/proto/ms"
	graphql1 "github.com/qilin/crm-api/generated/graphql"
)

type msMutationResolver struct{ *Resolver }
type msQueryResolver struct{ *Resolver }

func (r *Resolver) MsMutation() graphql1.MsMutationResolver {
	return &msMutationResolver{r}
}
func (r *Resolver) MsQuery() graphql1.MsQueryResolver {
	return &msQueryResolver{r}
}

func (r *msMutationResolver) New(ctx context.Context, obj *graphql1.MsMutation, name string) (*graphql1.NewOut, error) {
	c, d, e := getConnGRPC(context.Background(), srvMs)
	if e != nil {
		return nil, e
	}
	if d != nil {
		defer d()
	}
	cc := ms.NewMsClient(c)
	out, e := cc.New(ctx, &ms.NewIn{Name: name})
	if e != nil {
		return nil, e
	}
	return &graphql1.NewOut{
		Status: graphql1.NewOutStatus(out.Status.String()),
		ID:     typ.IntString(out.Id).V(),
	}, nil
}

func (r *msQueryResolver) Search(ctx context.Context, obj *graphql1.MsQuery, query string, cursor graphql1.CursorIn, order graphql1.OrderIn) (*graphql1.SearchOut, error) {
	c, d, e := getConnGRPC(context.Background(), srvMs)
	if e != nil {
		return nil, e
	}
	if d != nil {
		defer d()
	}
	cc := ms.NewMsClient(c)
	out, e := cc.Search(ctx, &ms.SearchIn{
		Query: query,
		Order: ms.Order(ms.Order_value[order.String()]),
		Cursor: &ms.CursorIn{
			Limit:  int64(cursor.Limit),
			Offset: int64(cursor.Offset),
			Cursor: typ.Of(cursor.Cursor).String().V(),
		},
	})
	if e != nil {
		return nil, e
	}
	ids := make([]string, len(out.Id))
	for i, id := range out.Id {
		ids[i] = typ.IntString(id).V()
	}
	if out.Cursor == nil {
		out.Cursor = &ms.CursorOut{}
	}
	return &graphql1.SearchOut{
		ID:     ids,
		Status: graphql1.SearchOutStatus(out.Status.String()),
		Cursor: &graphql1.CursorOut{
			Count:  int(out.Cursor.TotalCount),
			Limit:  int(out.Cursor.Limit),
			Offset: int(out.Cursor.Offset),
			IsEnd:  out.Cursor.HasNextPage,
			Cursor: out.Cursor.Cursor,
		},
	}, nil
}

func (r *mutationResolver) Ms(ctx context.Context) (*graphql1.MsMutation, error) {
	return &graphql1.MsMutation{}, nil
}

func (r *queryResolver) Ms(ctx context.Context) (*graphql1.MsQuery, error) {
	return &graphql1.MsQuery{}, nil
}
