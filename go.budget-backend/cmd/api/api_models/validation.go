package api_models

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ValidationError struct {
	Error     string `jsob:"error"`
	Key       string `json:"key"`
	Condition string `json:"condition"`
}

func ValidateBodyRequest(c echo.Context, payload interface{}) []*ValidationError {
	validate := validator.New(validator.WithRequiredStructEnabled())
	var errors []*ValidationError
	err := validate.Struct(payload)
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}
	for _, e := range validationErrors {
		c.Logger().Info(e.Field())
		fmt.Println(e.Field())
		currentError := &ValidationError{
			Error:     strings.Split(e.Error(), "Error:")[1],
			Key:       e.Field(),
			Condition: e.Tag(),
		}
		errors = append(errors, currentError)
	}

	return errors

}
