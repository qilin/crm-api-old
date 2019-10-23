package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/coreos/go-oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/qilin/crm-api/internal/db/domain"
	"golang.org/x/oauth2"
)

type Config struct {
	OAuth2 struct {
		Provider     string `required:true`
		ClientId     string `required:"true"`
		ClientSecret string `required:"true"`
		RedirectUrl  string `required:"true"`
	}
	AutoSignIn         bool
	Secret             string
	SuccessRedirectURL string
	JWT                struct {
		PublicKey  string
		PrivateKey string
	}
}

type Auth struct {
	cfg      *Config
	oauth2   *oauth2.Config
	verifier *oidc.IDTokenVerifier
	jwtKeys  KeyPair

	stateSecret []byte

	log logger.Logger

	users domain.UserRepo
}

func New(appCtx provider.LMT, cfg *Config, users domain.UserRepo) (*Auth, error) {
	var keys = oidc.NewRemoteKeySet(context.Background(), cfg.OAuth2.Provider+".well-known/jwks.json")

	jwtKeys, err := NewKeyPairFromPEM(cfg.JWT.PublicKey, cfg.JWT.PrivateKey)
	if err != nil {
		return nil, err
	}
	return &Auth{
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

		jwtKeys: jwtKeys,

		users: users,

		stateSecret: []byte(cfg.Secret),
		log:         appCtx.L().WithFields(logger.Fields{"service": "auth"}),
	}, nil
}

// Session ====================================================================

func (a *Auth) checkAuthorized(c echo.Context) (string, bool) {
	cssid, err := c.Cookie("ssid")
	if err == http.ErrNoCookie {
		return "", false
	}

	if err != nil {
		a.log.Warning("can't retrive cookies: %v", logger.Args(err))
		return "", false
	}

	a.log.Debug("session cookie: %s", logger.Args(cssid.Value))

	// validate jwt token
	if _, err := jwt.Parse(cssid.Value, func(*jwt.Token) (interface{}, error) {
		return a.jwtKeys.Public, nil
	}); err != nil {
		a.log.Debug("invalid session token: %v", logger.Args(err))
		return "", false
	}

	return cssid.Value, true
}

func (a *Auth) removeSession(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "ssid",
		Value:    "",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   false, // TODO
	})
}

func (a *Auth) setSession(c echo.Context, value string) {
	c.SetCookie(&http.Cookie{
		Name:     "ssid",
		Value:    value,
		MaxAge:   int((30 * time.Minute).Seconds()),
		HttpOnly: true,
		Secure:   false, // TODO
	})
}

// State ====================================================================

func (a *Auth) removeState(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "state",
		Value:    "",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   false, // TODO
	})
}

func (a *Auth) setState(c echo.Context, value string) {
	c.SetCookie(&http.Cookie{
		Name:     "state",
		Value:    value,
		MaxAge:   int((30 * time.Minute).Seconds()),
		HttpOnly: true,
		Secure:   false, // TODO
	})
}

func (a *Auth) validateState(c echo.Context, state string) error {
	cookie, err := c.Cookie("state")
	if err != nil {
		return err
	}

	if state != a.secureState(cookie.Value) {
		return errors.New("invalid auth state")
	}

	return nil
}

func (a *Auth) secureState(state string) string {
	var h = sha256.New()
	h.Write([]byte(state))
	h.Write(a.stateSecret)
	return hex.EncodeToString(h.Sum(nil))
}