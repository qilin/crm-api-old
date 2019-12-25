package providers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/google/wire"

	"github.com/ProtocolONE/go-core/v2/pkg/config"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/pkg/errors"

	"github.com/google/uuid"

	"github.com/qilin/crm-api/internal/authentication/common"

	"github.com/labstack/echo/v4"

	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

const (
	providerName = "p1"
)

type Config struct {
	Enabled bool
	Secret  string
	OAuth2  struct {
		Provider     string `required:"true"`
		ClientId     string `required:"true"`
		ClientSecret string `required:"true"`
		RedirectUrl  string `required:"true"`
	}
}

type P1Provider struct {
	ctx         context.Context
	cfg         *Config
	oauth2      *oauth2.Config
	verifier    *oidc.IDTokenVerifier
	stateSecret []byte
	provider.LMT
}

func NewP1Provider(ctx context.Context, set provider.AwareSet, cfg *Config) (*P1Provider, error) {
	keys := oidc.NewRemoteKeySet(context.Background(), cfg.OAuth2.Provider+".well-known/jwks.json")

	return &P1Provider{
		ctx: ctx,
		cfg: cfg,
		oauth2: &oauth2.Config{
			RedirectURL:  cfg.OAuth2.RedirectUrl,
			ClientID:     cfg.OAuth2.ClientId,
			ClientSecret: cfg.OAuth2.ClientSecret,
			Scopes:       []string{"openid"},
			Endpoint: oauth2.Endpoint{
				AuthURL:   cfg.OAuth2.Provider + "oauth2/auth",
				TokenURL:  cfg.OAuth2.Provider + "oauth2/token",
				AuthStyle: oauth2.AuthStyleInHeader,
			},
		},
		verifier: oidc.NewVerifier(cfg.OAuth2.Provider, keys, &oidc.Config{
			ClientID: cfg.OAuth2.ClientId,
		}),
		stateSecret: []byte(cfg.Secret),
		LMT:         &set,
	}, nil
}

func (a P1Provider) Provider() string {
	return providerName
}

func (a P1Provider) SignIn(ctx echo.Context) (user *common.ExternalUser, url string, err error) {
	var state = uuid.New().String()
	a.setState(ctx, state)
	return nil, a.oauth2.AuthCodeURL(a.secureState(state)), nil
}

func (a P1Provider) Callback(c echo.Context) (user *common.ExternalUser, err error) {
	id, err := a.callbackError(c)
	if err != nil {
		return nil, errors.New("failed")
	}

	return &common.ExternalUser{
		User:       common.User{},
		Provider:   "p1",
		ExternalId: id,
	}, nil
}

func (a *P1Provider) callbackError(c echo.Context) (string, error) {
	// Verify state param, defence from CSRF attacks
	var state = c.FormValue("state")
	a.L().Debug("oauth callback state %s", logger.Args(state))
	if err := a.validateState(c, state); err != nil {
		return "", err
	}
	a.removeState(c)

	// exchange code to tokens
	var code = c.FormValue("code")
	var ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	oauthToken, err := a.oauth2.Exchange(ctx, code)
	if err != nil {
		return "", errors.Wrap(err, "auth code exchange failed")
	}

	// parse and verify id_token
	rawIDToken, ok := oauthToken.Extra("id_token").(string)
	if !ok {
		return "", fmt.Errorf("id_token not provided")
	}

	idToken, err := a.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return "", errors.Wrap(err, "failed to verify id token")
	}

	a.L().Debug("user id: %s", logger.Args(idToken.Subject))

	return idToken.Subject, nil
}

// -- state ---
func (a *P1Provider) removeState(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "state",
		Value:    "",
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   false, // TODO
	})
}

func (a *P1Provider) setState(c echo.Context, value string) {
	c.SetCookie(&http.Cookie{
		Name:     "state",
		Value:    value,
		Path:     "/",
		MaxAge:   int((30 * time.Minute).Seconds()),
		HttpOnly: true,
		Secure:   false, // TODO
	})
}

func (a *P1Provider) validateState(c echo.Context, state string) error {
	cookie, err := c.Cookie("state")
	if err != nil {
		return err
	}

	if state != a.secureState(cookie.Value) {
		return errors.New("invalid auth state")
	}

	return nil
}

func (a *P1Provider) secureState(state string) string {
	var h = sha256.New()
	h.Write([]byte(state))
	h.Write(a.stateSecret)
	return hex.EncodeToString(h.Sum(nil))
}

func ProviderP1(ctx context.Context, set provider.AwareSet, cfg *Config) (*P1Provider, func(), error) {
	g, e := NewP1Provider(ctx, set, cfg)
	return g, func() {}, e
}

// Cfg
func Cfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{}
	e := cfg.UnmarshalKey("p1", c)
	return c, func() {}, e
}

// CfgTest
func CfgTest() (*Config, func(), error) {
	c := &Config{}
	return c, func() {}, nil
}

var (
	WireSet = wire.NewSet(
		ProviderP1,
		Cfg,
	)
)
