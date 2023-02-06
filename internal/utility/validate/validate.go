package validate

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Validate(v *validator.Validate, generic any) []string {
	err := v.Struct(generic)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return nil
		}

		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, fmt.Sprintf("%s is %s with type %s", err.StructField(), err.Tag(), err.Type()))
		}

		return errs
	}
	return nil
}
