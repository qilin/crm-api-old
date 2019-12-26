package main

import (
	"errors"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/qilin/crm-api/internal/authentication/common"
)

var (
	Plugin plugin
)

const (
	providerName = "rambler"
)

type plugin struct {
	//
}

func (p *plugin) SignIn(ctx echo.Context) (user *common.ExternalUser, url string, err error) {
	return &common.ExternalUser{
		User: common.User{
			FirstName: "John",
			LastName:  "Doe",
		},
		Provider:   providerName,
		ExternalId: time.Now().String(),
	}, "", nil
}

func (p *plugin) Callback(ctx echo.Context) (user *common.ExternalUser, err error) {
	return nil, errors.New("not implemented")
}

func (p *plugin) Provider() string {
	return providerName
}
