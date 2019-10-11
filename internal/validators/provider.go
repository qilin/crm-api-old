package validators

import (
	"github.com/google/wire"
	"gopkg.in/go-playground/validator.v9"
)

// Provider
func Provider() (*ValidatorSet, func(), error) {
	g := New()
	return g, func() {}, nil
}

// Validators
func ProviderValidators(v *ValidatorSet) (validate *validator.Validate, _ func(), err error) {
	validate = validator.New()

	if err = validate.RegisterValidation("password", v.PasswordValidator); err != nil {
		return
	}

	return validate, func() {}, nil
}

var (
	WireSet = wire.NewSet(
		Provider,
	)
	WireTestSet = wire.NewSet(
		Provider,
	)
)
