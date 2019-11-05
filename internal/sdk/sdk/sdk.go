package sdk

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"time"

	"github.com/qilin/crm-api/internal/db/domain"

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
	JWT     JWT
}

type JWT struct {
	Subject    string
	Iss        string
	Exp        int // time expiration in minutes
	PrivateKey string
	PublicKey  string
}

type KeyPair struct {
	Public  *ecdsa.PublicKey
	Private *ecdsa.PrivateKey
}

type SDK struct {
	ctx           context.Context
	cfg           Config
	pm            *plugins.PluginManager
	authenticator common.Authenticate
	orderer       common.Order
	keyRegister   jwt.KeyRegister
	repo          *repo.Repo
	keyPair       KeyPair
	provider.LMT
}

func (s *SDK) Mode() common.SDKMode {
	return s.cfg.Mode
}

func (s *SDK) Verify(token []byte) (*jwt.Claims, error) {
	// todo: return it back after tests
	// todo: optimise with ParseWithoutCheck + iss key map
	//return s.keyRegister.Check(token)
	return jwt.ECDSACheck(token, s.keyPair.Public)
}

func (s *SDK) Authenticate(ctx context.Context, request common.AuthRequest, token *jwt.Claims, log logger.Logger) (response common.AuthResponse, err error) {
	return s.authenticator(ctx, request, token, log)
}

func (s *SDK) Order(ctx context.Context, request common.OrderRequest, log logger.Logger) (response common.OrderResponse, err error) {
	return s.orderer(ctx, request, log)
}

func (s *SDK) GetProductByUUID(uuid string) (*domain.ProductItem, error) {
	return s.repo.Products.Get(s.ctx, uuid)
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

func New(ctx context.Context, set provider.AwareSet, repo *repo.Repo, cfg *Config) *SDK {
	pm := plugins.NewPluginManager()
	if cfg.Mode == common.ParentMode || cfg.Mode == common.DeveloperMode {
		for _, p := range cfg.Plugins {
			err := pm.Load(p)
			if err != nil {
				set.L().Emergency(err.Error())
			} else {
				set.L().Info("loaded plugin " + p)
			}
		}
	}

	sdk := &SDK{
		ctx:  ctx,
		cfg:  *cfg,
		repo: repo,
		pm:   pm,
		LMT:  &set,
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
		if err != nil {
			set.L().Emergency(err.Error())
		}

		qln := qilin.NewQilinAuthenticator()
		sdk.authenticator = qln.Auth
		sdk.orderer = qln.Order
	}
	return sdk
}

func decodePemECDSA(pemPriv, pemPub string) (KeyPair, error) {
	block, _ := pem.Decode([]byte(pemPriv))
	x509Encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		return KeyPair{}, err
	}

	blockPub, _ := pem.Decode([]byte(pemPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return KeyPair{}, err
	}
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return KeyPair{
		Public:  publicKey,
		Private: privateKey,
	}, nil
}
