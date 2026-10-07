package errors

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateObject[T any](data T) error {
	err := validate.Struct(data)

	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	var messages []string

	for _, fieldError := range validationErrors {
		field := fieldError.Field()

		switch fieldError.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s is required", field))

		case "max":
			messages = append(
				messages,
				fmt.Sprintf("%s must be at most %s characters", field, fieldError.Param()),
			)

		default:
			messages = append(
				messages,
				fmt.Sprintf("%s is invalid", field),
			)
		}
	}

	return fmt.Errorf(
		"missing or invalid field(s): %s",
		strings.Join(messages, "; "),
	)
}
