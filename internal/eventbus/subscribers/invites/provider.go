package invites

import (
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/google/wire"
)

func Cfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{
		invoker: invoker.NewInvoker(),
	}
	e := cfg.UnmarshalKeyOnReload(UnmarshalKey, c)
	return c, func() {}, e
}

// CfgTest
func CfgTest() (*Config, func(), error) {
	c := &Config{
		invoker: invoker.NewInvoker(),
	}
	return c, func() {}, nil
}

// Provider
func Provider(cfg *Config) (*InviteSubscriber, func(), error) {
	g := New(cfg)
	return g, func() {}, nil
}

var (
	WireSet = wire.NewSet(
		Provider,
		Cfg,
	)
	WireTestSet = wire.NewSet(
		Provider,
		CfgTest,
	)
)
