package handlers

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.budget-backend/cmd/api/api_models"
)

func (h *Handler) ValidateBodyRequest(c echo.Context, payload interface{}) []*api_models.ValidationError {
	validate := validator.New(validator.WithRequiredStructEnabled())
	var errors []*api_models.ValidationError
	err := validate.Struct(payload)
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}
	for _, e := range validationErrors {
		c.Logger().Info(e.Field())
		currentError := &api_models.ValidationError{
			Error:     strings.Split(e.Error(), "Error:")[1],
			Key:       e.Field(),
			Condition: e.Tag(),
		}
		errors = append(errors, currentError)
	}

	return errors

}
