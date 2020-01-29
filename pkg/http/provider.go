package http

import (
	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
)

// Cfg
func Cfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{}
	e := cfg.UnmarshalKey(UnmarshalKey, c)
	return c, func() {}, e
}

// CfgTest
func CfgTest() (*Config, func(), error) {
	c := &Config{}
	return c, func() {}, nil
}

// Provider
func Provider(set provider.AwareSet, dispatcher Dispatcher, cfg *Config) (*HTTP, func(), error) {
	http := New(set, dispatcher, cfg)
	return http, func() {}, nil
}

var (
	WireSet     = wire.NewSet(Provider, Cfg)
	WireTestSet = wire.NewSet(Provider, CfgTest)
)
