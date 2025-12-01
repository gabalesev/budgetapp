package api_models

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ApiResponse map[string]any

type JsonSuccessResponmse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type JsonFailedValidationResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Errors  []*ValidationError `json:"errors"`
}

type JsonErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func SendSuccessResponse(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, JsonSuccessResponmse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SendFailedValidationResponse(c echo.Context, message string, errors []*ValidationError) error {
	return c.JSON(http.StatusBadRequest, JsonFailedValidationResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

func SendErrorResponse(c echo.Context, message string, statusCode int) error {
	return c.JSON(statusCode, JsonErrorResponse{
		Success: false,
		Message: message,
	})
}

func SendBadRequestResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusBadRequest)
}

func SendNotFoundResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusNotFound)
}

func SendInternalServerErrorResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusInternalServerError)
}

func SendUnauthorizedResponse(c echo.Context, message string) error {

	if message == "" || len(message) == 0 {
		message = "Unauthorized"
	}
	return SendErrorResponse(c, message, http.StatusUnauthorized)
}
