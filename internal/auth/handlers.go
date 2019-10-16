package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var empty = map[string]interface{}{}

func (a *Auth) RegisterAPIGroup(ctx *echo.Echo) {
	var g = ctx.Group("/auth/v1")

	g.GET("/login", a.login)
	g.GET("/callback", a.callback)
	g.GET("/logout", a.logout)
	g.GET("/jwt", a.jwt)
}

func (a *Auth) login(c echo.Context) error {
	_, ok := a.checkAuthorized(c)
	if ok {
		return a.authSuccess(c)
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
	var state = c.FormValue("state")

	a.log.Debug("oauth callback state %s", logger.Args(state))

	if err := a.validateState(c, state); err != nil {
		return err
	}
	a.removeState(c)

	var code = c.FormValue("code")
	var ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	oauthToken, err := a.oauth2.Exchange(ctx, code)
	if err != nil {
		return errors.Wrap(err, "auth code exchange failed")
	}

	rawIDToken, ok := oauthToken.Extra("id_token").(string)
	if !ok {
		return fmt.Errorf("id_token not provided")
	}

	idtoken, err := a.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return errors.Wrap(err, "failed to verify id token")
	}

	a.log.Debug("user id: %s", logger.Args(idtoken.Subject))

	jwt, err := GenerateJWT(idtoken.Subject)
	if err != nil {
		return err
	}

	a.setSession(c, jwt)
	return a.authSuccess(c)
}

func (a *Auth) jwt(c echo.Context) error {
	if jwt, ok := a.checkAuthorized(c); ok {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"jwt": jwt,
		})
	}

	return c.JSON(http.StatusUnauthorized, empty)
}

func (a *Auth) authSuccess(c echo.Context) error {
	return c.Redirect(http.StatusFound, "http://localhost:3000/auth_success")
}
