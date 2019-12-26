package authentication

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
)

type Config struct {
	Debug                bool `fallback:"shared.debug"`
	Authenticator        string
	JWT                  JWT
	CookieName           string `default:"ssid"`
	CookieDomain         string
	CookieSecure         bool
	LoginSuccessRedirect string
	LogoutRedirect       string
	invoker              *invoker.Invoker
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
