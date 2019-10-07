package grpc

import (
	"context"
	"sync"
	"time"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"google.golang.org/grpc/keepalive"
)

var (
	pm   *PoolManager
	muPM sync.Mutex
)

// Cfg
func Cfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{
		invoker: invoker.NewInvoker(),
	}
	e := cfg.UnmarshalKeyOnReload(UnmarshalKey, c)
	return c, func() {}, e
}

// Service
type Service struct {
	Target          string
	MaxConn         int
	InitConn        int
	MaxLifeDuration time.Duration
	IdleTimeout     time.Duration
	Keepalive       *keepalive.ClientParameters
}

// Config
type Config struct {
	Services  map[string]*Service
	Keepalive *keepalive.ClientParameters
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

// Provider
func Provider(ctx context.Context, set provider.AwareSet, cfg *Config) (*PoolManager, func(), error) {
	muPM.Lock()
	defer muPM.Unlock()
	if pm != nil {
		return pm, func() {}, nil
	}
	pm = NewPoolManager(ctx, set, cfg)
	return pm, func() {}, nil
}

var WireSet = wire.NewSet(Provider, Cfg)
