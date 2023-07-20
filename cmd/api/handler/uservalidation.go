package handler

import (
	"strings"

	ut "github.com/go-playground/universal-translator"
	gpgvalidator "github.com/go-playground/validator/v10"
	"github.com/wisesight/go-api-template/pkg/validator"
)

const (
	isValidPasswordTag = "is_valid_password"
)

func isValidPassword(fl gpgvalidator.FieldLevel) bool {
	var (
		password = fl.Field().String()
		body     = fl.Parent().Interface().(CreateRequestBody)
		username = body.Username
	)
	return !strings.Contains(password, username)
}

func newUserValidation() error {
	errors := []error{
		// Validation
		validator.Validate.RegisterValidation(isValidPasswordTag, isValidPassword),

		// Translation
		validator.Validate.RegisterTranslation(
			isValidPasswordTag,
			validator.Trans,
			func(ut ut.Translator) error {
				return ut.Add(isValidPasswordTag, "password should not match the username", true)
			},
			func(ut ut.Translator, fe gpgvalidator.FieldError) string {
				t, _ := ut.T(isValidPasswordTag, fe.Field())
				return t
			},
		),
	}

	for _, err := range errors {
		if err != nil {
			// Should log here
			return err
		}
	}

	return nil
}
