package eventbus

import (
	"context"

	"github.com/qilin/crm-api/internal/stan"

	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
)

type EventBus struct {
	ctx  context.Context
	cfg  *Config
	stan *stan.Stan
	provider.LMT
}

// Config
type Config struct {
	Debug   bool `fallback:"shared.debug"`
	Secret  string
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

func New(ctx context.Context, set provider.AwareSet, stan *stan.Stan, cfg *Config) *EventBus {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": Prefix})
	return &EventBus{
		ctx:  ctx,
		cfg:  cfg,
		stan: stan,
		LMT:  &set,
	}
}
