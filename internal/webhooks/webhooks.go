package webhooks

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
	"github.com/qilin/crm-api/internal/dispatcher/common"
)

type WebHooks struct {
	ctx context.Context
	cfg *Config
	provider.LMT
}

// Route
func (h *WebHooks) Route(groups *common.Groups) {
	groups.V1.POST("/webhooks", h.handler)
}

func (h *WebHooks) handler(ctx echo.Context) error {
	if ctx.Request().Header.Get("X-Qilin-Secret") != h.cfg.Secret {
		return ctx.JSON(http.StatusForbidden, nil)
	}
	b, _ := ioutil.ReadAll(ctx.Request().Body)

	// process webhook
	// send it to service
	// go webhookProcessing()

	return ctx.JSON(http.StatusOK, map[string]string{"body": string(b), "secret": ctx.Request().Header.Get("X-Qilin-Secret")})
}

// Config
type Config struct {
	Debug   bool `fallback:"shared.debug"`
	Secret  string
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
func New(ctx context.Context, set provider.AwareSet, cfg *Config) *WebHooks {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": Prefix})
	return &WebHooks{
		ctx: ctx,
		cfg: cfg,
		LMT: &set,
	}
}
