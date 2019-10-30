package plugins

import "github.com/qilin/crm-api/internal/sdk/common"

type Authenticator interface {
	Name() string
	Auth(authenticate common.Authenticate) common.Authenticate
}

type Orderer interface {
	Name() string
	Order(order common.Order) common.Order
}
