package sdk

import (
	"context"
	"strings"

	"github.com/qilin/crm-api/internal/sdk/repo"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/pascaldekloe/jwt"
	"github.com/qilin/crm-api/internal/sdk/common"
	"github.com/qilin/crm-api/internal/sdk/plugins"
	"github.com/qilin/crm-api/internal/sdk/qilin"
)

type Config struct {
	Debug   bool `fallback:"shared.debug"`
	Mode    common.SDKMode
	Plugins []string
}

type SDK struct {
	ctx           context.Context
	cfg           Config
	pm            *plugins.PluginManager
	authenticator common.Authenticate
	orderer       common.Order
	keyRegister   jwt.KeyRegister
	provider.LMT
}

func (s *SDK) Mode() common.SDKMode {
	return s.cfg.Mode
}

func (s *SDK) Verify(token []byte) (*jwt.Claims, error) {
	return s.keyRegister.Check(token)
}

func (s *SDK) Authenticate(ctx context.Context, request common.AuthRequest, token *jwt.Claims, log logger.Logger) (response common.AuthResponse, err error) {
	return s.authenticator(ctx, request, token, log)
}

func (s *SDK) Order(ctx context.Context, request common.OrderRequest, log logger.Logger) (response common.OrderResponse, err error) {
	return s.orderer(ctx, request, log)
}

func New(ctx context.Context, set provider.AwareSet, repo *repo.Repo, cfg *Config) *SDK {
	pm := plugins.NewPluginManager()
	if cfg.Mode == common.ParentMode || cfg.Mode == common.DeveloperMode {
		for _, p := range cfg.Plugins {
			err := pm.Load(p)
			if err != nil {
				set.L().Error(err.Error())
			} else {
				set.L().Info("loaded plugin " + p)
			}
		}
	}

	sdk := &SDK{
		ctx: ctx,
		cfg: *cfg,
		pm:  pm,
		LMT: &set,
	}

	// load keys
	keys, err := repo.PlatformJWTKey.All(ctx)
	if err != nil {
		set.L().Emergency(err.Error())
	}
	for i := range keys {
		switch strings.ToLower(keys[i].KeyType) {
		case "pem":
			sdk.keyRegister.LoadPEM([]byte(keys[i].Key), []byte(""))
		case "jwk":
			sdk.keyRegister.LoadJWK([]byte(keys[i].Key))
		}
	}

	switch cfg.Mode {
	case common.ParentMode:
		// todo: configure parent sdk
		sdk.authenticator = pm.Auth(nil)
		break
	case common.DeveloperMode:
		// todo: configure developer sdk
		// todo: load plugins and build chain
		sdk.authenticator = pm.Auth(nil)
		break
	default:
		// todo: default qilin mode
		qln := qilin.NewQilinAuthenticator()
		sdk.authenticator = qln.Auth
		sdk.orderer = qln.Order
	}
	return sdk
}
