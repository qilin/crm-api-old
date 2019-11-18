package plugins

import (
	"context"
	"errors"
	"os"
	"plugin"

	"github.com/spf13/viper"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/qilin/crm-api/internal/sdk/common"
)

type PluginManager struct {
	init  []Initable
	auth  []Authenticator
	order []Orderer
	http  []Httper
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		init:  []Initable{},
		auth:  []Authenticator{},
		order: []Orderer{},
		http:  []Httper{},
	}
}

func (m *PluginManager) Load(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New("File not found {" + path + "}: " + err.Error())
	}

	file, err := plugin.Open(path)
	if err != nil {
		return errors.New("Can not open {" + path + "}: " + err.Error())
	}

	instance, err := file.Lookup("Plugin")
	if err != nil {
		return errors.New("Can not lookup {" + path + "}: " + err.Error())
	}

	init, ok := instance.(Initable)
	if ok {
		m.init = append(m.init, init)
	}

	ath, ok := instance.(Authenticator)
	if ok {
		m.auth = append([]Authenticator{ath}, m.auth...)
	}

	ord, ok := instance.(Orderer)
	if ok {
		m.order = append([]Orderer{ord}, m.order...)
	}

	httper, ok := instance.(Httper)
	if ok {
		m.http = append([]Httper{httper}, m.http...)
	}

	return nil
}

func (m *PluginManager) Init(ctx context.Context, cfg *viper.Viper, log logger.Logger) {
	for _, p := range m.init {
		p.Init(ctx, cfg, log)
	}
}

func (m *PluginManager) Http(ctx context.Context, echo2 *echo.Echo, log logger.Logger) {
	for _, p := range m.http {
		p.Http(ctx, echo2, log.WithFields(logger.Fields{"plugin": p.Name()}))
	}
}

func (m *PluginManager) Auth(authenticate common.Authenticate) common.Authenticate {
	for _, plg := range m.auth {
		authenticate = plg.Auth(authenticate)
	}
	return authenticate
}

func (m *PluginManager) Order(order common.Order) common.Order {
	for _, plg := range m.order {
		order = plg.Order(order)
	}
	return order
}
