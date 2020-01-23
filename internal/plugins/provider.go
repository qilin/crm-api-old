package plugins

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/google/wire"
)

func Cfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{}
	e := cfg.UnmarshalKeyOnReload(UnmarshalKey, c)
	return c, func() {}, e
}

// CfgTest
func CfgTest() (*Config, func(), error) {
	return &Config{}, func() {}, nil
}

// Provider
func Provider(ctx context.Context, cfg *Config, set provider.AwareSet, init config.Initial) (*PluginManager, func(), error) {
	pm := NewPluginManager(set.Logger)
	for _, path := range cfg.Plugins {
		err := pm.Load(path)
		if err != nil {
			set.L().Emergency(err.Error())
		} else {
			set.L().Info("loaded plugin from %s", logger.Args(path))
		}
	}
	sub := init.Viper.Sub(UnmarshalKeyPluginConfigs)
	if sub == nil {
		set.L().Warning("No plugins.configs defined")
		sub = init.Viper.Viper
	}
	pm.Init(ctx, sub, set.L())

	return pm, func() {}, nil
}

var (
	WireSet = wire.NewSet(
		Cfg,
		Provider,
	)

	WireTestSet = wire.NewSet(
		CfgTest,
		Provider,
	)
)
