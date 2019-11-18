package sdk

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path"
	"strings"
	"time"

	"github.com/ProtocolONE/go-core/v2/pkg/config"

	"github.com/davecgh/go-spew/spew"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
	"github.com/qilin/crm-api/internal/db/domain"
	"github.com/qilin/crm-api/internal/sdk/common"
	"github.com/qilin/crm-api/internal/sdk/plugins"
	"github.com/qilin/crm-api/internal/sdk/qilin"
	"github.com/qilin/crm-api/internal/sdk/repo"
)

type interfaces struct {
	authenticator common.Authenticate
	orderer       common.Order
}

type SDK struct {
	ctx           context.Context
	cfg           *Config
	pluginsCfg    *PluginsConfig
	pluginManager *plugins.PluginManager
	interfaces    interfaces
	keyRegister   jwt.KeyRegister
	repo          *repo.Repo
	keyPair       KeyPair
	provider.LMT
}

func (s *SDK) Authenticate(ctx context.Context, request common.AuthRequest, token *jwt.Claims, log logger.Logger) (response common.AuthResponse, err error) {
	spew.Dump(s.pluginsCfg.PluginsConfig)
	return s.interfaces.authenticator(context.WithValue(ctx, "config", s.pluginsCfg.PluginsConfig), request, token, log)
}

func (s *SDK) GetProductByUUID(uuid string) (*domain.ProductItem, error) {
	//return s.repo.Products.Get(s.ctx, uuid)
	url, ok := s.cfg.Iframes[uuid]
	if ok {
		// todo: hardcoded reponse with config
		return &domain.ProductItem{
			ID:        uuid,
			URL:       url, // todo: return config value
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}
	return nil, fmt.Errorf("unknown product UUID '%s'", uuid)
}

func (s *SDK) IframeHtml(qiliProductUUID string) (string, error) {
	// in qilin mode no difference, always same html for all products
	tplName := path.Base(s.cfg.IframeTemplate)
	tpl, err := template.New(tplName).ParseFiles(s.cfg.IframeTemplate)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, nil)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *SDK) IssueJWT(userId, qilinProductUUID string) ([]byte, error) {
	claims := jwt.Claims{
		Set: map[string]interface{}{},
	}
	claims.Subject = s.cfg.JWT.Subject
	claims.Issuer = s.cfg.JWT.Iss
	claims.Set[common.UserID] = userId
	claims.Set[common.QilinProductUUID] = qilinProductUUID

	if s.cfg.JWT.Exp > 0 {
		now := time.Now().Round(time.Second)
		claims.Expires = jwt.NewNumericTime(now.Add(time.Duration(s.cfg.JWT.Exp) * time.Minute))
	}

	// issue a JWT
	return claims.ECDSASign(jwt.ES512, s.keyPair.Private)
}

func (s *SDK) MapExternalUserToUser(platformId int, externalId string) string {
	//user, err := s.repo.UserMap.FindByExternalID(s.ctx, platformId, externalId)
	// if user not found, create new user, otherwise return user.UserId
	//if err == gorm.ErrRecordNotFound {
	//	user := &domain.UserMapItem{
	//		PlatformID: platformId,
	//		ExternalID: externalId,
	//	}
	//	s.repo.UserMap.Create(s.ctx, user)
	//}
	return "1c3e43a5-8513-42b3-8774-596c78079bb2"
}

func (s *SDK) Mode() common.SDKMode {
	return s.cfg.Mode
}

func (s *SDK) Order(ctx context.Context, request common.OrderRequest, log logger.Logger) (response common.OrderResponse, err error) {
	return s.interfaces.orderer(context.WithValue(ctx, "config", s.pluginsCfg.PluginsConfig), request, log)
}

func (s *SDK) PluginsRoute(echo *echo.Echo) {
	s.pluginManager.Http(context.WithValue(context.Background(), "config", s.pluginsCfg.PluginsConfig), echo, s.L())
}

func (s *SDK) Verify(token []byte) (*jwt.Claims, error) {
	// todo: optimise with ParseWithoutCheck + iss key map
	//return s.keyRegister.Check(token) // todo: uncomment it back after tests
	return jwt.ECDSACheck(token, s.keyPair.Public)
}

func New(ctx context.Context, set provider.AwareSet, repo *repo.Repo, cfg *Config, pCfg *PluginsConfig, init config.Initial) *SDK {
	pm := plugins.NewPluginManager()
	if cfg.Mode == common.StoreMode || cfg.Mode == common.ProviderMode {
		for _, p := range cfg.Plugins {
			err := pm.Load(p)
			if err != nil {
				set.L().Emergency(err.Error())
			} else {
				set.L().Info("loaded plugin " + p)
			}
		}
	}

	pm.Init(ctx, init.Viper.Sub(common.UnmarshalKeyPluginConfig), set.L())

	sdk := &SDK{
		ctx:           ctx,
		cfg:           cfg,
		pluginsCfg:    pCfg,
		repo:          repo,
		pluginManager: pm,
		LMT:           &set,
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

	sdk.keyPair, err = decodePemECDSA(cfg.JWT.PrivateKey, cfg.JWT.PublicKey)
	if err != nil {
		set.L().Emergency(err.Error())
	}

	switch cfg.Mode {
	case common.StoreMode:
		sdk.interfaces.authenticator = pm.Auth(nil)
		break
	case common.ProviderMode:
		sdk.interfaces.authenticator = pm.Auth(nil)
		break
	default:
		qln := qilin.NewQilinAuthenticator()
		sdk.interfaces.authenticator = qln.Auth
		sdk.interfaces.orderer = qln.Order
	}
	return sdk
}
