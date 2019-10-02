package daemon

import (
	"context"
	"net"

	"github.com/qilin/crm-api/generated/api/proto/ms"
	"github.com/qilin/crm-api/internal/db/domain"
	"github.com/qilin/crm-api/internal/db/repo"
	"github.com/qilin/crm-api/internal/db/trx"
	"github.com/qilin/go-core/invoker"
	"github.com/qilin/go-core/logger"
	"github.com/qilin/go-core/provider"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// Daemon
type Daemon struct {
	ctx    context.Context
	cfg    Config
	server *grpc.Server
	repo   Repo
	trx    *trx.Manager
	provider.LMT
}

func (d *Daemon) Search(ctx context.Context, in *ms.SearchIn) (*ms.SearchOut, error) {
	out := &ms.SearchOut{Status: ms.SearchOut_OK}
	//
	cursor := &domain.Cursor{}
	cursor.Limit.Set(int(in.Cursor.Limit))
	cursor.Offset.Set(int(in.Cursor.Offset))
	cursor.Cursor.Set(in.Cursor.Cursor)
	//
	l, e := d.repo.List.List(ctx, []string{"id"}, cursor, domain.ToOrder(in.Order.String()), in.Query)
	if e != nil {
		switch {
		case repo.IsRecordNotFoundError(e):
			out.Status = ms.SearchOut_NOT_FOUND
		default:
			out.Status = ms.SearchOut_SERVER_INTERNAL_ERROR
		}
		return out, nil
	}
	//
	out.Id = make([]int64, len(l))
	for i, item := range l {
		out.Id[i] = item.ID.Typ().Int64().V()
	}
	out.Cursor = &ms.CursorOut{
		TotalCount:  int64(cursor.TotalCount.V()),
		Limit:       int64(cursor.Limit.V()),
		Offset:      int64(cursor.Offset.V()),
		HasNextPage: cursor.HasNextPage.V(),
		Cursor:      cursor.Cursor.V(),
	}
	return out, nil
}

func (d *Daemon) New(ctx context.Context, in *ms.NewIn) (*ms.NewOut, error) {
	out := &ms.NewOut{Status: ms.NewOut_OK}
	//
	li := &domain.ListItem{}
	li.Name.Set(in.Name)
	//
	e := d.repo.List.Create(ctx, li)
	if e != nil {
		out.Status = ms.NewOut_SERVER_INTERNAL_ERROR
		return out, nil
	}
	out.Id = li.ID.Typ().Int64().V()
	return out, nil
}

// ListenAndServe
func (d *Daemon) ListenAndServe() (err error) {

	ms.RegisterMsServer(d.server, d)

	d.L().Info("start listen and serve at %v", logger.Args(d.cfg.Bind))

	go func() {
		<-d.ctx.Done()
		d.L().Info("context cancelled, shutdown is raised")
		d.server.GracefulStop()
	}()

	listener, err := net.Listen("tcp", d.cfg.Bind)
	if err != nil {
		d.L().Error("server is shutdown with error, %v", logger.Args(err))
	}
	return d.server.Serve(listener)
}

// Keepalive defines configurable parameters for point-to-point healthcheck.
type Keepalive struct {
	ServerParameters  keepalive.ServerParameters
	EnforcementPolicy keepalive.EnforcementPolicy
}

// Config is a general GRPC config settings
type Config struct {
	Debug     bool   `fallback:"shared.debug"`
	Bind      string `required:"true"`
	Keepalive Keepalive
	invoker   *invoker.Invoker
}

// OnReload
func (c *Config) OnReload(callback func(ctx context.Context)) {
	c.invoker.OnReload(callback)
}

// Reload
func (c *Config) Reload(ctx context.Context) {
	c.invoker.Reload(ctx)
}

type AppSet struct {
	Repo Repo
	Trx  *trx.Manager
}

// New
func New(ctx context.Context, set provider.AwareSet, appSet AppSet, cfg *Config) *Daemon {
	s := grpc.NewServer(
		grpc.KeepaliveParams(cfg.Keepalive.ServerParameters),
		grpc.KeepaliveEnforcementPolicy(cfg.Keepalive.EnforcementPolicy),
	)
	set.Logger = set.Logger.WithFields(logger.Fields{"service": Prefix})
	return &Daemon{
		server: s,
		ctx:    ctx,
		cfg:    *cfg,
		repo:   appSet.Repo,
		trx:    appSet.Trx,
		LMT:    &set,
	}
}
