package grpc

import (
	"context"

	"github.com/qilin/go-core/logger"
	"github.com/qilin/go-core/provider"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// PoolManager
type PoolManager struct {
	ctx context.Context
	cfg Config
	provider.LMT
}

// New
func (p *PoolManager) New(service string) (_ *Pool, loaded bool, _ error) {
	s, ok := p.cfg.Services[service]
	if !ok {
		return nil, false, errCfgInvalid
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	cl := s.Keepalive
	if cl == nil {
		if p.cfg.Keepalive != nil {
			cl = p.cfg.Keepalive
		} else {
			cl = &keepalive.ClientParameters{}
		}
	}
	opts = append(opts, grpc.WithKeepaliveParams(*cl))
	pool, l := NewPool(p.ctx, service, s.Target,
		MaxConn(s.MaxConn),
		InitConn(s.InitConn),
		MaxLifeDuration(s.MaxLifeDuration),
		IdleTimeout(s.IdleTimeout),
		ConnOptions(opts...),
	)
	return pool, l, nil
}

// NewPoolManager
func NewPoolManager(ctx context.Context, set provider.AwareSet, cfg *Config) *PoolManager {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": Prefix})
	return &PoolManager{
		ctx: ctx,
		cfg: *cfg,
		LMT: &set,
	}
}
