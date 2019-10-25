package auth

import (
	"context"
	"fmt"
	"github.com/qilin/crm-api/internal/dispatcher/common"
	"net/http"
	"time"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/qilin/crm-api/internal/db/domain"
)

var empty = map[string]interface{}{}

func (a *Auth) Route(groups *common.Groups) {
	groups.Auth.GET("/login", a.login)
	groups.Auth.GET("/callback", a.callback)
	groups.Auth.GET("/logout", a.logout)
	groups.Auth.GET("/jwt", a.jwt)
}

func (a *Auth) login(c echo.Context) error {
	_, ok := a.checkAuthorized(c)
	if ok {
		return c.Redirect(http.StatusFound, a.cfg.SuccessRedirectURL)
	}

	var state = uuid.New().String()
	a.setState(c, state)

	var url = a.oauth2.AuthCodeURL(a.secureState(state))
	return c.Redirect(http.StatusFound, url)
}

func (a *Auth) logout(c echo.Context) error {
	a.removeSession(c)
	return c.JSON(http.StatusOK, empty)
}

func (a *Auth) callback(c echo.Context) error {

	// Verify state param, defence from CSRF attacks
	var state = c.FormValue("state")
	a.L().Debug("oauth callback state %s", logger.Args(state))
	if err := a.validateState(c, state); err != nil {
		return err
	}
	a.removeState(c)

	// exchange code to tokens
	var code = c.FormValue("code")
	var ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	oauthToken, err := a.oauth2.Exchange(ctx, code)
	if err != nil {
		return errors.Wrap(err, "auth code exchange failed")
	}

	// parse and verify id_token
	rawIDToken, ok := oauthToken.Extra("id_token").(string)
	if !ok {
		return fmt.Errorf("id_token not provided")
	}

	idToken, err := a.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return errors.Wrap(err, "failed to verify id token")
	}

	a.L().Debug("user id: %s", logger.Args(idToken.Subject))

	// find user
	u, err := a.appSet.UserRepo.FindByExternalID(context.TODO(), idToken.Subject)
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			a.L().Debug("%v", logger.Args(err))
			return err
		}

		if !a.cfg.AutoSignIn {
			a.L().Debug("user, not registered")
			return errors.New("not registered")
		}

		var claims struct {
			Email         string `json:"email"`
			EmailVerified bool   `json:"email_verified"`
		}
		if err := idToken.Claims(&claims); err != nil {
			return err
		}

		a.L().Debug("user not found, create new one")
		if err := a.appSet.UserRepo.Create(ctx, &domain.UserItem{
			TenantID:   1, // TODO
			ExternalID: idToken.Subject,
			Email:      claims.Email,
			Role:       "owner", // TODO default role ?
		}); err != nil {
			a.L().Debug("%v", logger.Args(err))
			return err
		}
		u, err = a.appSet.UserRepo.FindByExternalID(ctx, idToken.Subject)
		if err != nil {
			a.L().Debug("%v", logger.Args(err))
			return err
		}
	}

	a.L().Info("user logged in %d", logger.Args(u.ID))

	// create auth jwt
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, NewClaims(u))
	signed, err := token.SignedString(a.jwtKeys.Private)
	if err != nil {
		return err
	}

	a.setSession(c, signed)
	return c.Redirect(http.StatusFound, a.cfg.SuccessRedirectURL)
}

func (a *Auth) jwt(c echo.Context) error {
	if token, ok := a.checkAuthorized(c); ok {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"jwt": token,
		})
	}
	return c.JSON(http.StatusUnauthorized, empty)
}
