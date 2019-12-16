package cache

import (
	"context"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/cache"
	"github.com/eko/gocache/store"
	"github.com/pkg/errors"
)

// Cache
type Cache struct {
	ctx   context.Context
	cfg   *Config
	store store.StoreInterface
	provider.LMT
}

// Get
func (c Cache) Get(key interface{}) (interface{}, error) {
	return c.store.Get(key)
}

// Set
func (c Cache) Set(key, object interface{}, options *store.Options) error {
	return c.store.Set(key, object, options)
}

// Delete
func (c Cache) Delete(key interface{}) error {
	return c.store.Delete(key)
}

// Invalidate
func (c Cache) Invalidate(options store.InvalidateOptions) error {
	return c.store.Invalidate(options)
}

// Clear
func (c Cache) Clear() error {
	return c.store.Clear()
}

// GetType
func (c Cache) GetType() string {
	return c.store.GetType()
}

// Config
type Config struct {
	Debug     bool   `fallback:"shared.debug"`
	Store     string `default:"ristretto"`
	Ristretto *ristretto.Config
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

// New
func New(ctx context.Context, set provider.AwareSet, cfg *Config) (*Cache, error) {
	var m store.StoreInterface
	switch cfg.Store {
	case "ristretto":
		c, err := ristretto.NewCache(cfg.Ristretto)
		if err != nil {
			return nil, errors.WithMessage(err, Prefix)
		}
		s := store.NewRistretto(c, nil)
		m = cache.New(s)
	default:
		return nil, ErrNotSupported
	}
	set.Logger = set.Logger.WithFields(logger.Fields{"service": Prefix})
	return &Cache{
		store: m,
		ctx:   ctx,
		cfg:   cfg,
		LMT:   &set,
	}, nil
}
