package validators

import (
	"gopkg.in/go-playground/validator.v9"
)

type ValidatorSet struct {
	validate *validator.Validate
}

// Email validator
func (v *ValidatorSet) PasswordValidator(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) > 4
}

func New() *ValidatorSet {
	return &ValidatorSet{}
}
