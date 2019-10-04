package jwt

import (
	"github.com/google/wire"
	"github.com/qilin/go-core/config"
	"github.com/qilin/go-core/invoker"
)

func Provider(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{
		invoker: invoker.NewInvoker(),
	}
	e := cfg.UnmarshalKeyOnReload(UnmarshalKey, c)
	return c, func() {}, e
}

var (
	WireSet = wire.NewSet(
		Provider,
	)
	WireTestSet = wire.NewSet(
		Provider,
	)
)
