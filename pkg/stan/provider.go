package stan

import (
	"context"

	"github.com/google/wire"
	"github.com/qilin/go-core/config"
	"github.com/qilin/go-core/provider"
)

// Cfg
func Cfg(cfg config.Configurator) (Config, func(), error) {
	c := Config{}
	e := cfg.UnmarshalKey(UnmarshalKey, &c)
	return c, func() {}, e
}

// CfgTest
func CfgTest() (Config, func(), error) {
	return Config{}, func() {}, nil
}

// Provider
func Provider(ctx context.Context, set provider.AwareSet, cfg *Config) (*Stan, func(), error) {
	g := New(ctx, set, cfg)
	return g, func() {}, nil
}

var (
	WireSet     = wire.NewSet(Provider, Cfg)
	WireTestSet = wire.NewSet(Provider, CfgTest)
)
