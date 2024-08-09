package vutils

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)

	val := reflect.TypeOf(s)

	// Check if the input is a pointer and dereference it
	if val.Kind() != reflect.Pointer {
		panic("ValidateStruct only accepts pointer")
	}

	for _, err := range err.(validator.ValidationErrors) {
		// Reflectively get the field value
		field, _ := val.Elem().FieldByName(err.Field())
		fieldName := field.Name
		var message string
		message, ok := field.Tag.Lookup("error_message")

		if !ok {
			switch err.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", fieldName)
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters long", fieldName, err.Param())
			case "max":
				message = fmt.Sprintf("%s must be at most %s characters long", fieldName, err.Param())
			case "email":
				message = fmt.Sprintf("Invalid email format")
			case "gte":
				message = fmt.Sprintf("%s must be greater than or equal to %s", fieldName, err.Param())
			case "lte":
				message = fmt.Sprintf("%s must be less than or equal to %s", fieldName, err.Param())
			default:
				message = fmt.Sprintf("%s is not valid", fieldName)
			}
		}

		errors[fieldName] = message
	}

	return errors
}
