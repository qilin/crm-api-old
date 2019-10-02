package postgres

import (
	"context"
	"time"

	"github.com/qilin/go-core/invoker"
)

const (
	Prefix       = "go-shared.postgres"
	UnmarshalKey = "postgres"
)

// Config is a general db config settings
type Config struct {
	Debug           bool `fallback:"shared.debug"`
	Dsn             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	invoker         *invoker.Invoker
}

// OnReload
func (c *Config) OnReload(callback func(ctx context.Context)) {
	c.invoker.OnReload(callback)
}

// Reload
func (c *Config) Reload(ctx context.Context) {
	c.invoker.Reload(ctx)
}
