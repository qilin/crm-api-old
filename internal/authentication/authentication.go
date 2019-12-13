package authentication

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/labstack/echo/v4"
	"github.com/qilin/crm-api/internal/auth"
	"github.com/qilin/crm-api/internal/crypto"
	"github.com/qilin/crm-api/internal/db/domain"
	"github.com/qilin/crm-api/internal/plugins"
)

type (
	Authentication interface {
		// check
		IsAuthenticated(ctx context.Context) bool
		Middleware(next echo.HandlerFunc) echo.HandlerFunc
		// session
		SignIn(ctx context.Context) (string, error)
		SignOut(ctx context.Context) error
		SetSession(c echo.Context, value string)
		RemoveSession(c echo.Context)
	}

	Authenticator interface {
		SignIn(ctx echo.Context) (*User, *Redirect, error)
		SignOut(ctx echo.Context) error
		Callback(ctx echo.Context) (r *Redirect, err error)
	}

	Redirect struct {
	}

	AuthenticationService struct {
		auth       Authenticator
		app        AppSet
		cfg        Config
		jwtKeyPair crypto.KeyPair
		provider.LMT
	}
)

func (a *AuthenticationService) SignIn(ctx context.Context) (string, error) {
	user := ExtractUserContext(ctx)
	if user == nil {
		a.L().Error("user is not authenticated")
		return "", NotAuthenticated{}
	}

	claims := &SessionClaims{
		User: User{
			ID:       user.ID,
			Language: user.Language,
		},
		StandardClaims: jwt.StandardClaims{
			Audience: "", // todo
			IssuedAt: time.Now().Unix(),
			Issuer:   "", // todo
			Subject:  "", // todo
		},
	}
	token, err := a.jwtKeyPair.Sign(claims)
	if err != nil {
		a.L().Error("SignIn error: ", logger.Args(err))
		return "", err
	}

	err = a.app.AuthLog.Create(ctx, &domain.AuthLogItem{
		UserID:    user.ID,
		Action:    domain.ActionSignIn,
		UserAgent: "", // todo: extract from context
		IP:        "", // todo: extract from context
		HWID:      "", // todo: extract from context
		CreatedAt: time.Now(),
	})
	if err != nil {
		a.L().Error("SignIn error: ", logger.Args(err))
		return "", err
	}

	return token, nil
}

func (a *AuthenticationService) SignOut(ctx context.Context) error {
	user := ExtractUserContext(ctx)
	if user == nil {
		return NotAuthenticated{}
	}
	err := a.app.AuthLog.Create(ctx, &domain.AuthLogItem{
		UserID:    user.ID,
		Action:    domain.ActionSignOut,
		UserAgent: "", // todo: extract from context
		IP:        "", // todo: extract from context
		HWID:      "", // todo: extract from context
		CreatedAt: time.Now(),
	})
	if err != nil {
		a.L().Error("SignOut error: ", logger.Args(err))
		return err
	}

	return nil
}

func (a *AuthenticationService) IsAuthenticated(ctx context.Context) bool {
	return IsAuthenticated(ctx)
}

func (a *AuthenticationService) SetSession(c echo.Context, value string) {
	c.SetCookie(&http.Cookie{
		Name:     a.cfg.CookieName,
		Value:    value,
		HttpOnly: true,
		Domain:   a.cfg.CookieDomain,
		Path:     "/",
		Secure:   a.cfg.CookieSecure,
	})
}

func (a *AuthenticationService) RemoveSession(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     a.cfg.CookieName,
		Value:    "",
		HttpOnly: true,
		Domain:   a.cfg.CookieDomain,
		Path:     "/",
		Secure:   a.cfg.CookieSecure,
	})
}

// Middleware extracts session info from cookie into context
func (a *AuthenticationService) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessCookie, err := ctx.Request().Cookie(a.cfg.CookieName)
		if err != nil {
			//a.L().Error("no session cookie")
			return next(ctx)
		}

		// validate jwt token
		var claims = &SessionClaims{}
		if err := a.jwtKeyPair.Parse(sessCookie.Value, claims); err != nil {
			a.L().Error("invalid jwt token: %v", logger.Args(err))
			return next(ctx)
		}

		SetUserContext(ctx, &User{
			ID:       claims.User.ID,
			Language: claims.User.Language,
		})
		a.L().Debug("auth user: %d", logger.Args(auth.ExtractUserContext(ctx.Request().Context())))

		return next(ctx)
	}
}

type AppSet struct {
	AuthLog         domain.AuthLogRepo
	UsersRepo       domain.UsersRepo
	UserProviderMap domain.UserProviderMapRepo
}

func New(ctx context.Context, pm *plugins.PluginManager, set provider.AwareSet, app AppSet, cfg *Config) (*AuthenticationService, error) {
	// todo: select authenticator from config
	kp, err := crypto.NewKeyPairFromPEM(cfg.JWT.PublicKey, cfg.JWT.PrivateKey)
	set.L().Emergency("Can not parse authentication JWT key pair: %s", logger.Args(err))

	// todo: load auth from plugin or user default
	auth := &AuthenticationService{
		jwtKeyPair: kp,
		auth:       nil,
		LMT:        &set,
	}

	return auth, nil
}
