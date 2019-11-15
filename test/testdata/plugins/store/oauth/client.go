package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ProtocolONE/go-core/v2/pkg/logger"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

var RamblerEndpoint = oauth2.Endpoint{
	AuthURL:   "https://id.rambler.ru/oauthd/",
	TokenURL:  "https://id.rambler.ru/oauthsrv/access_token",
	AuthStyle: oauth2.AuthStyleInParams,
}

const (
	getUserDataURL = "https://id.rambler.ru/oauthsrv/api/getCurrentUser?access_token="
)

type Auth struct {
	Log          logger.Logger
	RedirectURL  string
	ClientID     string
	ClientSecret string
	Scopes       []string
	Endpoint     oauth2.Endpoint
}

type RamblerUserInfo struct {
	Id        string
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Status    string
	Email     string
	Gender    string
	Birthday  string
}

func (a *Auth) AuthHandler(ctx echo.Context) error {
	oauthConfig := a.ramblerOauthConfig()
	return ctx.Redirect(http.StatusTemporaryRedirect, oauthConfig.AuthCodeURL(""))
}

func (a *Auth) AuthCallbackHandler(ctx echo.Context) error {
	content, err := a.getUserInfo(ctx.Request().FormValue("state"), ctx.Request().FormValue("code"))
	if err != nil {
		a.Log.Error(err.Error())
		http.Redirect(ctx.Response().Writer, ctx.Request(), "/", http.StatusTemporaryRedirect)
		return err
	}

	user := RamblerUserInfo{}
	err = json.Unmarshal(content, &user)
	if err != nil {
		a.Log.Error(err.Error())
		return err
	}

	// todo: fill cookie with JWT
	cookie := &http.Cookie{}
	cookie.Value = "jwt-token"
	ctx.SetCookie(cookie)
	return nil
}

func (a *Auth) getUserInfo(state string, code string) ([]byte, error) {
	oauthConfig := a.ramblerOauthConfig()
	if state != "pseudo-random" {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get(getUserDataURL + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}

// https://github.com/douglasmakey/oauth2-example
func (a *Auth) ramblerOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  a.RedirectURL,
		ClientID:     a.ClientID,
		ClientSecret: a.ClientSecret,
		Scopes:       []string{},
		Endpoint:     RamblerEndpoint,
	}
}
