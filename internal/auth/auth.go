package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/coreos/go-oidc"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	qilinCrypto "github.com/qilin/crm-api/internal/crypto"
	"github.com/qilin/crm-api/internal/db/domain"
	"github.com/qilin/crm-api/internal/dispatcher/common"
	"golang.org/x/oauth2"
)

type Config struct {
	Enabled bool
	OAuth2  struct {
		Provider     string `required:"true"`
		ClientId     string `required:"true"`
		ClientSecret string `required:"true"`
		RedirectUrl  string `required:"true"`
	}
	LogoutRedirect string

	// cookies rules
	SessionCookieName string `default:"ssid"`
	SessionTTL        int    `default:336` // session max age in hours
	Domain            string
	SecureCookie      bool

	AutoSignIn         bool
	Secret             string
	SuccessRedirectURL string
	ErrorRedirectURL   string
	JWT                struct {
		PublicKey  string
		PrivateKey string
	}
}

type Auth struct {
	cfg         Config
	oauth2      *oauth2.Config
	verifier    *oidc.IDTokenVerifier
	jwtKeys     qilinCrypto.KeyPair
	stateSecret []byte
	appSet      AppSet
	provider.LMT
}

type AppSet struct {
	UserRepo domain.UserRepo
}

// New
func New(set provider.AwareSet, appSet AppSet, cfg *Config) (*Auth, error) {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": common.Prefix})
	keys := oidc.NewRemoteKeySet(context.Background(), cfg.OAuth2.Provider+".well-known/jwks.json")

	var jwtKeys qilinCrypto.KeyPair
	if cfg.Enabled {
		k, err := qilinCrypto.NewKeyPairFromPEM(cfg.JWT.PublicKey, cfg.JWT.PrivateKey)
		if err != nil {
			return nil, errors.Wrap(err, "can't parse auth.jwt keys")
		}
		jwtKeys = k
	}

	return &Auth{
		cfg: *cfg,
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
		jwtKeys:     jwtKeys,
		stateSecret: []byte(cfg.Secret),
		appSet:      appSet,
		LMT:         &set,
	}, nil
}

// Session ====================================================================

func (a *Auth) checkAuthorized(c echo.Context) (*AccessTokenClaims, bool) {
	var token string
	auth := c.Request().Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		token = auth[7:]
	} else {
		ssid, err := c.Cookie(a.cfg.SessionCookieName)
		if err == http.ErrNoCookie {
			return nil, false
		}

		if err != nil {
			a.L().Warning("can't retrieve cookies: %v", logger.Args(err))
			return nil, false
		}
		token = ssid.Value
	}

	a.L().Debug("auth token: %s", logger.Args(token))

	// validate jwt token
	var claims AccessTokenClaims
	if err := a.jwtKeys.Parse(token, &claims); err != nil {
		a.L().Debug("invalid auth token: %v", logger.Args(err))
		return nil, false
	}

	if claims.UserID != "" {
		id, err := strconv.Atoi(claims.UserID)
		if err != nil {
			a.L().Warning("%v", logger.Args(err))
			return nil, false
		}
		u, err := a.appSet.UserRepo.Get(c.Request().Context(), id)
		if err != nil {
			a.L().Warning("%v", logger.Args(err))
			return nil, false
		}

		if u.AuthTimestamp.Unix() > claims.IssuedAt {
			a.L().Debug("auth token revoked")
			return nil, false
		}

	}

	return &claims, true
}

func (a *Auth) removeSession(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     a.cfg.SessionCookieName,
		Value:    "",
		HttpOnly: true,
		Domain:   a.cfg.Domain,
		Path:     "/",
		Secure:   a.cfg.SecureCookie,
	})
}

func (a *Auth) setSession(c echo.Context, value string) {
	c.SetCookie(&http.Cookie{
		Name:     a.cfg.SessionCookieName,
		Value:    value,
		HttpOnly: true,
		MaxAge:   int((time.Duration(a.cfg.SessionTTL) * time.Hour).Seconds()),
		Domain:   a.cfg.Domain,
		Path:     "/",
		Secure:   a.cfg.SecureCookie,
		SameSite: http.SameSiteLaxMode,
	})
}

// State ====================================================================

func (a *Auth) removeState(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "state",
		Value:    "",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   a.cfg.SecureCookie,
	})
}

func (a *Auth) setState(c echo.Context, value string) {
	c.SetCookie(&http.Cookie{
		Name:     "state",
		Value:    value,
		Domain:   a.cfg.Domain,
		MaxAge:   int((30 * time.Minute).Seconds()),
		HttpOnly: true,
		Secure:   a.cfg.SecureCookie,
		SameSite: http.SameSiteNoneMode,
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

// redirectURL ==============================================================

const (
	redirectURLCookie = "redirect"
)

// todo: make it safe (store on backend or encrypt / sign)
func (a *Auth) saveRedirectURL(c echo.Context, url string) {
	// todo: probably it would be better to save it to redis or so
	c.Set(ctxRedirectURL, url)
	c.SetCookie(&http.Cookie{
		Name:     redirectURLCookie,
		Value:    url,
		Domain:   a.cfg.Domain,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   a.cfg.SecureCookie,
		SameSite: http.SameSiteNoneMode,
	})
}

func (a *Auth) extractRedirectURL(c echo.Context) string {
	var url string
	cookie, err := c.Cookie(redirectURLCookie)
	if err != nil {
		a.L().Error(err.Error())
		return ""
	}
	if cookie.Value != "" {
		url = cookie.Value
	}
	if url == "" {
		url, _ = c.Get(ctxRedirectURL).(string)
	}
	return url
}

func (a *Auth) removeRedirectURL(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     redirectURLCookie,
		Value:    "",
		Domain:   a.cfg.Domain,
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   a.cfg.SecureCookie,
	})
}
