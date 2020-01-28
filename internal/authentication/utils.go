package authentication

import (
	"github.com/qilin/crm-api/internal/authentication/common"
	"github.com/qilin/crm-api/internal/db/domain"
)

func item2user(user *domain.UsersItem) *common.User {
	return &common.User{
		ID:    user.ID,
		Email: user.Email,
		//Phone:     user.Phone, // todo: tmp fix
		//Address1:  user.Address1,
		//Address2:  user.Address2,
		//City:      user.City,
		//State:     user.State,
		//Country:   user.Country,
		//Zip:       user.Zip,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		//BirthDate: user.BirthDate, // todo: tmp fix
		//Language:  user.Language,
	}
}
