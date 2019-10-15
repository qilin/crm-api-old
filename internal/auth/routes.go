package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/oauth2"
)

type Auth struct {
	issuer string
	oauth2 *oauth2.Config
	keys   jwt.KeyRegister
}

func New() *Auth {
	var keys jwt.KeyRegister
	n, err := keys.LoadJWK([]byte(`
{
    "keys": [
        {
            "use": "sig",
            "kty": "RSA",
            "kid": "public:e196b1ae-d0f2-46c8-b4f5-57d5861ee226",
            "alg": "RS256",
            "n": "uqiGwjdFRyGvUWFwsNBO-kRy5oGA-KN_gs1P6cLpcegt--3-V9kIY4QhYt5-0xAWbESI6d1keQDPU3-71zoFoKIW2NSBxnsrcMKOhNgo46matWWpHP52aZilyyhIPahMLmr2alVZemqKSy1yk5ZGSz7BvuBijMmAPh3q1NFAweS9wHKLdEfEaRJCX13zMFdNvh70J8PRwwwgG0gUNV5_10TliKPReg9hQ0amzaJBjDSxBaMP-Ai0DqJQhBNjHh4BIeDCk2XzouyF1BrvbkBGu9T-g8XSxkRxVk-k9U_2x0tc2eyD3h8BKi-OpIBpE6vStfPwhYcmRkqayC-GEyqB-_FMZHEbYKnm19nmC0pq3TtvxwqUJcJemDsxtZSSJcGKoX-u1Rot6S5KnclTq2F3h9b76f7tcCbgc8_OOJsQmu8rKYgaXv6Fd4nWzj4Alw7B2CGkPxO3wGtUK1XB0JOKk7HmfbcfhW5VMX0Bz7Id2Y0fuCbGeMaEq2sXGmH4UugTAZYUZ01YW8acjgALRjef7tLVhTPrSQ7d4ciZ3ShZfy1Xbur_CgltDdau3lRmWGWyfHe9-FXY6WREfpNDGP5FIdOh6fE7JRqW6c1_f18wjBTILXFbzURhTX86LMvTG1bIuS5lDKufce8Gk3zGOCeg1TcHN7lRYTxuu7lfAynId80",
            "e": "AQAB"
        }
    ]
}
`))

	if err != nil {
		panic(err)
	}

	return &Auth{
		issuer: "https://auth1.tst.protocol.one/",
		keys:   keys,
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
	}
}

func (a *Auth) RegisterAPIGroup(ctx *echo.Echo) {
	var g = ctx.Group("/auth/v1")

	g.GET("/login", a.login)
	g.GET("/callback", a.callback)
	g.GET("/logout", a.logout)
}

func (a *Auth) login(c echo.Context) error {
	cssid, err := c.Cookie("ssid")
	if err == nil {
		fmt.Println("cookie:", cssid.Value)
		if ValidateJWT(cssid.Value) {
			fmt.Println("cookie valid")
			q := url.Values{}
			q.Add("auth_token", cssid.Value)

			return c.Redirect(http.StatusFound, "http://localhost:8080?"+q.Encode())
		}
		fmt.Println("cookie not valid")
	}

	var url = a.oauth2.AuthCodeURL("some long state")
	return c.Redirect(http.StatusFound, url)
}

func (a *Auth) logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "ssid",
		Value:    "",
		Expires:  time.Now().Add(-30 * time.Minute),
		HttpOnly: true,
		Secure:   false, // TODO
	})

	q := url.Values{}
	q.Add("auth_token", "")

	return c.Redirect(http.StatusFound, "http://localhost:8080?"+q.Encode())
}

func (a *Auth) callback(c echo.Context) error {
	var (
		state = c.FormValue("state")
		code  = c.FormValue("code")
	)
	fmt.Println(state) // TODO state

	token, err := a.oauth2.Exchange(oauth2.NoContext, code) //TODO context
	if err != nil {
		return fmt.Errorf("code exchange failed: %s", err.Error())
	}

	fmt.Println(token)

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return fmt.Errorf("id_token not provided")
	}
	_ = rawIDToken

	claims, err := a.keys.Check([]byte(rawIDToken))
	if err != nil {
		return err
	}

	if !claims.AcceptAudience(a.oauth2.ClientID) {
		return fmt.Errorf("wrong OpenID audience")
	}

	if claims.Issuer != a.issuer {
		return fmt.Errorf("wrong OpenID issuer")
	}

	if claims.Expires.Time().Before(time.Now()) {
		return fmt.Errorf("OpenID token expired")
	}

	fmt.Println("User ID:", claims.Subject)

	jwt, err := GenerateJWT(claims.Subject)

	c.SetCookie(&http.Cookie{
		Name:     "ssid",
		Value:    jwt,
		Expires:  time.Now().Add(30 * time.Minute),
		HttpOnly: true,
		Secure:   false, // TODO
	})

	q := url.Values{}
	q.Add("auth_token", jwt)

	return c.Redirect(http.StatusFound, "http://localhost:8080?"+q.Encode())
}
