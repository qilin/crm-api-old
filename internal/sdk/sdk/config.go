package sdk

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/qilin/crm-api/internal/sdk/common"
)

type Config struct {
	Debug   bool `fallback:"shared.debug"`
	Mode    common.SDKMode
	Iframes map[string]string // todo: it's temporary
	Plugins []string
	JWT     JWT
	invoker *invoker.Invoker
}

type JWT struct {
	Subject    string
	Iss        string
	Exp        int // time expiration in minutes
	PrivateKey string
	PublicKey  string
}

// OnReload
func (c *Config) OnReload(callback func(ctx context.Context)) {
	c.invoker.OnReload(callback)
}

// Reload
func (c *Config) Reload(ctx context.Context) {
	c.invoker.Reload(ctx)
}
