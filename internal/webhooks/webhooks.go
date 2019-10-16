package webhooks

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/qilin/crm-api/internal/eventbus/events"

	"github.com/gurukami/typ/v2"

	"github.com/qilin/crm-api/internal/eventbus"

	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
	"github.com/qilin/crm-api/internal/dispatcher/common"
)

type WebHooks struct {
	ctx context.Context
	cfg *Config
	eb  *eventbus.EventBus
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

	h.L().Info(string(b))

	var hook Hook
	e := json.Unmarshal(b, &hook)
	if e != nil {
		h.L().Error("unmarshal failed with error: %v", logger.Args(e.Error()))
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	exp, _ := time.Parse("2000-01-31T11:00:00+00:00", typ.Of(hook.Event.Data.New["expiration"]).String().V())
	inv := events.Invite{
		Id:         typ.Of(hook.Event.Data.New["id"]).Int(0).V(),
		TenantId:   typ.Of(hook.Event.Data.New["tenant_id"]).Int(0).V(),
		UserId:     typ.Of(hook.Event.Data.New["user_id"]).Int(0).V(),
		Expiration: exp,
		Email:      typ.Of(hook.Event.Data.New["email"]).String().V(),
		FirstName:  typ.Of(hook.Event.Data.New["first_name"]).String().V(),
		LastName:   typ.Of(hook.Event.Data.New["last_name"]).String().V(),
		Accepted:   false,
	}
	go func() {
		e := h.eb.Publish(inv)
		if e != nil {
			h.L().Error(e.Error())
		}
	}()

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
func New(ctx context.Context, set provider.AwareSet, eb *eventbus.EventBus, cfg *Config) *WebHooks {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": Prefix})
	return &WebHooks{
		ctx: ctx,
		cfg: cfg,
		eb:  eb,
		LMT: &set,
	}
}
