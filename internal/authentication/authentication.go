package authentication

import (
	"context"
	"net/http"
	"time"

	"github.com/qilin/crm-api/internal/authentication/providers"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/ProtocolONE/go-core/v2/pkg/provider"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	common2 "github.com/qilin/crm-api/internal/authentication/common"
	"github.com/qilin/crm-api/internal/crypto"
	"github.com/qilin/crm-api/internal/db/domain"
	"github.com/qilin/crm-api/internal/dispatcher/common"
	"github.com/qilin/crm-api/internal/plugins"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

type (
	AuthenticationService struct {
		providers  map[string]common2.AuthenticationProvider
		app        AppSet
		cfg        *Config
		validate   *validator.Validate
		jwtKeyPair crypto.ECDSAKeyPair
		provider.LMT
	}
)

func (a *AuthenticationService) SignUp(ctx echo.Context) error {
	r := new(SignUpRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"status": "400",
		})
	}

	exists, err := a.app.UsersRepo.IsExistsEmail(ctx.Request().Context(), r.Email)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"status": "500",
		})
	}
	if exists {
		return ctx.JSON(http.StatusConflict, map[string]string{
			"status":  "409",
			"message": "email already in use",
		})
	}

	// validate
	if err := a.validate.Var(r.Email, "email,required"); err != nil {
		return ctx.JSON(http.StatusConflict, map[string]string{
			"status":  "400",
			"message": err.Error(),
		})
	}
	if err := a.validate.Var(r.Password, "password,required"); err != nil {
		return ctx.JSON(http.StatusConflict, map[string]string{
			"status":  "400",
			"message": err.Error(),
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		a.L().Error("SignUp error: ", logger.Args(err))
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"status": "500",
		})
	}

	user := &domain.UsersItem{
		Email:    r.Email,
		Password: string(hash),
		Status:   domain.UserActive.Int8(),
		// todo: fill other fields?
		CreatedAt: time.Now(),
	}
	err = a.app.UsersRepo.Create(ctx.Request().Context(), user)
	if err != nil {
		a.L().Error("SignUp error: ", logger.Args(err))
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"status": "500",
		})
	}

	err = a.app.AuthLog.Create(ctx.Request().Context(), &domain.AuthLogItem{
		UserID:    user.ID,
		Action:    domain.ActionSignUp,
		UserAgent: ctx.Request().UserAgent(),
		IP:        ctx.Request().RemoteAddr,
		HWID:      "", // todo: extract from request / context
		CreatedAt: time.Now(),
	})
	if err != nil {
		a.L().Error("SignUp error: ", logger.Args(err))
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"status": "500",
		})
	}

	// user registered
	return ctx.JSON(http.StatusOK, map[string]string{
		"status": "200",
	})
}

func (a *AuthenticationService) SignIn(ctx echo.Context) error {
	if a.IsAuthenticated(ctx.Request().Context()) {
		a.L().Debug("externalUser already authenticated")
		return ctx.JSON(http.StatusOK, map[string]string{
			"status": "200",
		})
	}

	providerName := ctx.Param("provider")
	if providerName == "" {
		// login by email & password
		a.L().Info("log by name / password")
		return a.loginByEmailPassword(ctx)
	}

	provider := a.getProvider(providerName)
	if provider == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"status":  "404",
			"message": "provider not found",
		})
	}

	externalUser, url, err := provider.SignIn(ctx)
	// error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"status": "500",
		})
	}
	// redirect
	if url != "" {
		return ctx.Redirect(http.StatusTemporaryRedirect, url)
	}
	// no redirect & empty external user
	if externalUser == nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"status": "500",
		})
	}

	// map external user to user
	user, err := a.mapUser(ctx.Request().Context(), *externalUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"status": "500",
		})
	}

	err = a.startSession(ctx, *user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"status": "500",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"status": "200",
	})
}

func (a *AuthenticationService) Callback(ctx echo.Context) error {
	provider := ctx.Param("provider")
	if provider == "" {
		provider = "p1"
	}

	externalUser, err := a.getProvider(provider).Callback(ctx)
	if err != nil {
		a.L().Error("callback error %s", logger.Args(err.Error()))
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"status": "500",
		})
	}

	user, err := a.mapUser(ctx.Request().Context(), *externalUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"status": "500",
		})
	}

	// start session
	if err := a.startSession(ctx, *user); err != nil {
		a.L().Error("Authentication.Callback error on startSession: %s", logger.Args(err))
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"status": "500",
		})
	}

	return ctx.Redirect(http.StatusTemporaryRedirect, a.cfg.LoginSuccessRedirect)
	//return ctx.JSON(http.StatusOK, map[string]string{
	//	"status": "200",
	//})
}

func (a *AuthenticationService) SignOut(ctx echo.Context) error {
	if !a.IsAuthenticated(ctx.Request().Context()) {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"status": "401",
		})
	}

	// remove session
	a.removeSession(ctx)

	user := ExtractUserContext(ctx.Request().Context())
	err := a.app.AuthLog.Create(ctx.Request().Context(), &domain.AuthLogItem{
		UserID:    user.ID,
		Action:    domain.ActionSignOut,
		UserAgent: ctx.Request().UserAgent(),
		IP:        ctx.Request().RemoteAddr,
		HWID:      "", // todo: extract from context
		CreatedAt: time.Now(),
	})
	if err != nil {
		a.L().Error("SignOut error: ", logger.Args(err))
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"status": "500",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"status": "200",
	})
}

// Middleware extracts session info from cookie into context
func (a *AuthenticationService) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessCookie, err := ctx.Request().Cookie(a.cfg.CookieName)
		if err != nil {
			a.L().Debug("no session cookie")
			return next(ctx)
		}

		// validate jwt token
		var claims = &SessionClaims{}
		if err := a.jwtKeyPair.Parse(sessCookie.Value, claims); err != nil {
			a.L().Error("invalid jwt token: %v", logger.Args(err))
			return next(ctx)
		}
		user := &common2.User{
			ID:       claims.User.ID,
			Language: claims.User.Language,
		}
		SetUserContext(ctx, user)
		a.L().Debug("auth user: %d", logger.Args(ExtractUserContext(ctx.Request().Context())))

		return next(ctx)
	}
}

func (a *AuthenticationService) IsAuthenticated(ctx context.Context) bool {
	return IsAuthenticated(ctx)
}

func (a *AuthenticationService) startSession(ctx echo.Context, user common2.User) error {
	claims := &SessionClaims{
		User: common2.User{
			ID:       user.ID,
			Language: user.Language,
		},
		StandardClaims: jwt.StandardClaims{
			Audience: "store", // todo
			IssuedAt: time.Now().Unix(),
			Issuer:   "iss.store", // todo
			Subject:  "session",   // todo
		},
	}
	token, err := a.jwtKeyPair.Sign(claims)
	if err != nil {
		a.L().Error("KEY SIGN ERROR")
		return err
	}

	err = a.app.AuthLog.Create(ctx.Request().Context(), &domain.AuthLogItem{
		UserID:    user.ID,
		Action:    domain.ActionSignIn,
		UserAgent: ctx.Request().UserAgent(),
		IP:        ctx.Request().RemoteAddr,
		HWID:      "", // todo: extract from context
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	// set cookie
	a.setSession(ctx, token)

	return nil
}

func (a *AuthenticationService) setSession(c echo.Context, value string) {
	c.SetCookie(&http.Cookie{
		Name:     a.cfg.CookieName,
		Value:    value,
		HttpOnly: false,
		Domain:   a.cfg.CookieDomain,
		Path:     "/",
		Secure:   a.cfg.CookieSecure,
	})
}

func (a *AuthenticationService) removeSession(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     a.cfg.CookieName,
		Value:    "",
		HttpOnly: true,
		Domain:   a.cfg.CookieDomain,
		Path:     "/",
		Secure:   a.cfg.CookieSecure,
	})
}

func (a *AuthenticationService) getProvider(name string) common2.AuthenticationProvider {
	p, ok := a.providers[name]
	if !ok {
		return nil
	}
	return p
}

func (a *AuthenticationService) mapUser(ctx context.Context, externalUser common2.ExternalUser) (*common2.User, error) {
	e2u, err := a.app.UserProviderMap.Get(ctx, externalUser.Provider, externalUser.ExternalId)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		a.L().Error("Authentication.Callback error on UserProviderMap.Get(%s,%s): %s", logger.Args(externalUser.Provider, externalUser.ExternalId, err))
		return nil, err
	}

	var u *common2.User
	// map external externalUser to externalUser if mapping not found
	if gorm.IsRecordNotFoundError(err) {
		a.L().Info("CREATE NEW EXTERNAL USER")
		user := &domain.UsersItem{
			Email:        externalUser.Email,
			Phone:        externalUser.Phone,
			Password:     "",
			Status:       domain.UserActive.Int8(),
			ServiceLevel: 0,
			Address1:     externalUser.Address1,
			Address2:     externalUser.Address2,
			City:         externalUser.City,
			State:        externalUser.State,
			Country:      externalUser.Country,
			Zip:          externalUser.Zip,
			PhotoURL:     "", // todo
			FirstName:    externalUser.FirstName,
			LastName:     externalUser.LastName,
			BirthDate:    externalUser.BirthDate,
			Language:     externalUser.Language,
			CreatedAt:    time.Now(),
		}
		// todo: run in transaction
		err := a.app.UsersRepo.Create(ctx, user)
		if err != nil {
			//
		}
		err = a.app.UserProviderMap.Create(ctx, &domain.UserProviderMapItem{
			UserID:      externalUser.ID,
			Provider:    externalUser.Provider,
			ProviderID:  externalUser.ExternalId,
			ProviderKey: "",
			CreatedAt:   time.Now(),
		})
		u = item2user(user)
	} else {
		user, err := a.app.UsersRepo.Get(ctx, e2u.UserID)
		if err != nil {
			a.L().Error("Authentication.Callback error on UsersRepo.Get(%d): %s", logger.Args(e2u.UserID, err))
			return nil, err
		}
		u = item2user(user)
	}

	return u, nil
}

func (a *AuthenticationService) loginByEmailPassword(ctx echo.Context) error {
	r := new(SignInRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"status": "400",
		})
	}

	// todo: validate email & password

	user, err := a.app.UsersRepo.FindByEmail(ctx.Request().Context(), r.Email)
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			a.L().Error("Authentication.loginByEmailPassword with email (%s) error: %s", logger.Args(r.Email, err))
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"status": "500",
			})
		}
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"status": "400",
		})
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Password)) != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"status": "400",
		})
	}

	u := item2user(user)
	err = a.startSession(ctx, *u)

	if err != nil {
		a.L().Error("Authentication.startSession with email (%s) error: %s", logger.Args(r.Email, err))
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"status": "500",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"status": "200",
	})
}

// Routes
func (a *AuthenticationService) Route(groups *common.Groups) {
	groups.Auth.POST("/signup", a.SignUp)
	groups.Auth.POST("/login", a.SignIn)
	groups.Auth.GET("/login/:provider", a.SignIn)
	groups.Auth.GET("/login/:provider/callback", a.Callback)
	groups.Auth.GET("/logout", a.SignOut)
}

type AppSet struct {
	AuthLog         domain.AuthLogRepo
	UsersRepo       domain.UsersRepo
	UserProviderMap domain.UserProviderMapRepo
}

func New(ctx context.Context, pm *plugins.PluginManager, set provider.AwareSet, app AppSet, validate *validator.Validate, p1 *providers.P1Provider, cfg *Config) (*AuthenticationService, error) {
	kp, err := crypto.NewECDSAKeyPairFromPEM(cfg.JWT.PublicKey, cfg.JWT.PrivateKey)
	if err != nil {
		set.L().Emergency("Can not parse authentication JWT key pair: %s", logger.Args(err))
	}

	// todo: load auth from plugin or user default
	auth := &AuthenticationService{
		providers: map[string]common2.AuthenticationProvider{
			p1.Provider(): p1,
		},
		cfg:        cfg,
		jwtKeyPair: kp,
		app:        app,
		validate:   validate,
		LMT:        &set,
	}

	for _, ap := range pm.AuthProviders() {
		auth.providers[ap.Provider()] = ap
	}

	return auth, nil
}
