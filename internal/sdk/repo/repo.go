package repo

import (
	"github.com/qilin/crm-api/internal/db/domain"
)

type Repo struct {
	Store       domain.StoreRepo
	StoreJWTKey domain.StoreJWTKeyRepo
	StoreGames  domain.StoreGamesRepo
	UserMap     domain.UserMapRepo
	ActionsLog  domain.ActionsLog
}
