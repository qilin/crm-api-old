package validators

import (
	"gopkg.in/go-playground/validator.v9"
)

type ValidatorSet struct {
	validate *validator.Validate
}

// Email validator
func (v *ValidatorSet) EmailValidator(email string) bool {
	return false
}

func New() *ValidatorSet {
	return &ValidatorSet{}
}
