package auth

import (
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

// Cfg
func Cfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{}
	err := cfg.UnmarshalKey(UnmarshalKey, c)
	return c, func() {}, errors.Wrap(err, "failed to parse auth configuration")
}

// CfgTest
func CfgTest() (*Config, func(), error) {
	c := &Config{}
	return c, func() {}, nil
}

// Provider
func Provider(set provider.AwareSet, appSet AppSet, cfg *Config) (*Auth, func(), error) {
	g, e := New(set, appSet, cfg)
	return g, func() {}, e
}

var (
	WireSet     = wire.NewSet(Provider, Cfg, wire.Struct(new(AppSet), "*"))
	WireTestSet = wire.NewSet(Provider, CfgTest, wire.Struct(new(AppSet), "*"))
)
