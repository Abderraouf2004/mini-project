package apperr

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	err := validate.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		return strings.TrimSpace(fl.Field().String()) != ""
	})
	if err != nil {
		panic(err)
	}
}

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
		field := jsonFieldName(data, fieldError.Field())

		switch fieldError.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s is required", field))

		case "notblank":
			messages = append(messages, fmt.Sprintf("%s cannot be blank", field))

		case "max":
			messages = append(messages, fmt.Sprintf(
				"%s must be at most %s characters",
				field,
				fieldError.Param(),
			))

		case "oneof":
			messages = append(messages, fmt.Sprintf(
				"%s has an invalid value",
				field,
			))

		default:
			messages = append(messages, fmt.Sprintf("%s is invalid", field))
		}
	}

	return errors.New(strings.Join(messages, "; "))
}

func jsonFieldName(data any, fieldName string) string {
	t := reflect.TypeOf(data)

	if t == nil {
		return fieldName
	}

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return fieldName
	}

	field, ok := t.FieldByName(fieldName)
	if !ok {
		return fieldName
	}

	name := strings.Split(field.Tag.Get("json"), ",")[0]
	if name == "" || name == "-" {
		return fieldName
	}

	return name
}

func ValidateAndConvert[T any](data any) (T, error) {
	var result T

	raw, err := json.Marshal(data)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(raw, &result); err != nil {
		return result, err
	}

	if err := ValidateObject(result); err != nil {
		return result, fmt.Errorf("response DTO validation: %w", err)
	}

	return result, nil
}
