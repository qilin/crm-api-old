package auth

import (
	"context"
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
)

// Cfg
func Cfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{
	}
	e := cfg.UnmarshalKey(UnmarshalKey, c)
	return c, func() {}, e
}

// CfgTest
func CfgTest() (*Config, func(), error) {
	c := &Config{
	}
	return c, func() {}, nil
}

// Provider
func Provider(ctx context.Context, set provider.AwareSet, appSet AppSet, cfg *Config) (*Auth, func(), error) {
	g, e := New(ctx, set, appSet, cfg)
	return g, func() {}, e
}

var (
	WireSet     = wire.NewSet(Provider, Cfg, wire.Struct(new(AppSet), "*"), )
	WireTestSet = wire.NewSet(Provider, CfgTest, wire.Struct(new(AppSet), "*"), )
)
