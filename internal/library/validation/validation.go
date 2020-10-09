package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func New() *validator.Validate {
	return validator.New()
}

// Validate takes in a model (book, author) and uses a validator to validate
func Validate(val *validator.Validate, model interface{}) []string {
	err := val.Struct(model)
	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("the field \"%s\" is %s", err.Field(), err.Tag())
			errs = append(errs, errorMessage)
		}

		return errs
	}
	return nil
}
