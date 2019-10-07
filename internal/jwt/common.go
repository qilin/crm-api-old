package jwt

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
)

const (
	Prefix       = "internal.jwt"
	UnmarshalKey = "dispatcher.jwt"
)

// Config is a general db config settings
type Config struct {
	Alg     string
	Private string
	Public  string
	invoker *invoker.Invoker
}

// OnReload
func (c *Config) OnReload(callback func(ctx context.Context)) {
	c.invoker.OnReload(callback)
}

// Reload
func (c *Config) Reload(ctx context.Context) {
	c.invoker.Reload(ctx)
}
