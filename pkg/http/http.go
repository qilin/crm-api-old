package http

import (
	"context"
	"net/http"

	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
)

// HTTP
type HTTP struct {
	cfg        Config
	dispatcher Dispatcher
	provider.LMT
}

// ListenAndServe
func (h *HTTP) ListenAndServe(ctx context.Context) (err error) {

	server := echo.New()
	server.HideBanner = true
	server.HidePort = true
	server.Debug = h.cfg.Debug

	if err := h.dispatcher.Dispatch(server); err != nil {
		return err
	}

	h.L().Info("start listen and serve http at %v", logger.Args(h.cfg.Bind))

	go func() {
		<-ctx.Done()
		h.L().Info("context cancelled, shutdown is raised")
		if e := server.Shutdown(context.Background()); e != nil {
			h.L().Error("graceful shutdown error, %v", logger.Args(e))
		}
	}()

	if err = server.Start(h.cfg.Bind); err != nil {
		if err == http.ErrServerClosed {
			err = nil
		} else {
			return err
		}
	}

	h.L().Info("http server stopped successfully")
	return nil
}

// Config
type Config struct {
	Debug   bool   `fallback:"shared.debug"`
	Bind    string `required:"true"`
	invoker *invoker.Invoker
}

// OnReload
func (c *Config) OnReload(callback func(ctx context.Context)) {
	c.invoker.OnReload(callback)
}

// Reload
func (c *Config) Reload(ctx context.Context) {
	c.invoker.Reload(ctx)
}

// New
func New(set provider.AwareSet, dispatcher Dispatcher, cfg *Config) *HTTP {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": Prefix})
	return &HTTP{
		cfg:        *cfg,
		dispatcher: dispatcher,
		LMT:        &set,
	}
}
