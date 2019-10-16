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
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type Auth struct {
	issuer   string
	oauth2   *oauth2.Config
	log      logger.Logger
	verifier *oidc.IDTokenVerifier

	stateSecret []byte
}

func New(appCtx provider.LMT) *Auth {
	var issuer = "https://auth1.tst.protocol.one/"
	keys := oidc.NewRemoteKeySet(context.Background(), issuer+".well-known/jwks.json")

	return &Auth{
		issuer: issuer,
		oauth2: &oauth2.Config{
			RedirectURL:  "http://localhost:8082/auth/v1/callback",
			ClientID:     "5da4ec412b13220001efe179",
			ClientSecret: "w9JkFlYOa7Hr6QY4OwMv4CWUl3rkJJkwkBjEdkGFE6SQILaM3Hpzxsvs5REdMFNV",
			Scopes:       []string{"openid"},
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://auth1.tst.protocol.one/oauth2/auth",
				TokenURL:  "https://auth1.tst.protocol.one/oauth2/token",
				AuthStyle: oauth2.AuthStyleInHeader,
			},
		},

		log: appCtx.L().WithFields(logger.Fields{"service": "auth"}),
		verifier: oidc.NewVerifier(issuer, keys, &oidc.Config{
			ClientID: "5da4ec412b13220001efe179",
		}),

		stateSecret: []byte("some secret"), // TODO config
	}
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

	if !ValidateJWT(cssid.Value) {
		a.log.Debug("session token expired or invalid")
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
