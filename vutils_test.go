package vutils

import (
	"fmt"
	"testing"
)

type User struct {
	Name  string `validate:"required,min=3,max=10"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=18,lte=60" error_message:"My custom error message"`
}

func TestValidateStruct(t *testing.T) {
	t.Run("Valid fields", func(t *testing.T) {
		user := &User{
			Name:  "John Doe",
			Email: "john@doe.com",
			Age:   18,
		}

		err := ValidateStruct(user)

		if err != nil {
			t.Errorf("Incorrect result. Expected nil, got %v", err)
		}
	})

	t.Run("Invalid field", func(t *testing.T) {
		user := &User{
			Name:  "J",
			Email: "john@doe.com",
			Age:   18,
		}

		err := ValidateStruct(user)

		if err == nil {
			t.Errorf("Incorrect result. Expected map[string]string, got %v", err)
		}

		err_length := len(err)
		if err_length != 1 {
			t.Errorf("Incorrect result. Expected length of map err to be 1, got %v", err_length)
		}
		expected_error := fmt.Sprintf("%s must be at least %v characters long", "Name", 3)

		if err["Name"] != expected_error {
			t.Errorf("Incorrect result. Expected error message to be \"%v\" got \"%v\"", expected_error, err["Name"])
		}
	})

	t.Run("Invalid field with custom error messages", func(t *testing.T) {
		user := &User{
			Name:  "John Doe",
			Email: "john@doe.com",
			Age:   17,
		}

		err := ValidateStruct(user)

		if err == nil {
			t.Errorf("Incorrect result. Expected map[string]string, got %v", err)
		}

		err_length := len(err)
		if err_length != 1 {
			t.Errorf("Incorrect result. Expected length of map err to be 1, got %v", err_length)
		}
		expected_error := "My custom error message"

		if err["Age"] != expected_error {
			t.Errorf("Incorrect result. Expected error message to be \"%v\" got \"%v\"", expected_error, err["Age"])
		}
	})

	t.Run("Uncovered validation tag", func(t *testing.T) {
		type User struct {
			Name  string `validate:"required,min=3,max=10"`
			Email string `validate:"required,email"`
			Age   int    `validate:"gte=18,lte=60" error_message:"My custom error message"`
			IP    string `validate:"cidr"`
		}

		user := &User{
			Name:  "John Doe",
			Email: "john@doe.com",
			Age:   18,
			IP:    "100000",
		}

		err := ValidateStruct(user)

		if err == nil {
			t.Errorf("Incorrect result. Expected map[string]string, got %v", err)
		}

		err_length := len(err)
		if err_length != 1 {
			t.Errorf("Incorrect result. Expected length of map err to be 1, got %v", err_length)
		}
		expected_error := fmt.Sprintf("%s is not valid", "IP")

		if err["IP"] != expected_error {
			t.Errorf("Incorrect result. Expected error message to be \"%v\" got \"%v\"", expected_error, err["IP"])
		}
	})
}
