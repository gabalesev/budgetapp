package handlers

import (
	"encoding/base64"
	"errors"

	"net/url"

	"github.com/labstack/echo/v4"
	"go.budget-backend/cmd/api/api_models"
	"go.budget-backend/cmd/api/services"
	"go.budget-backend/internal/mailer"
	"gorm.io/gorm"
)

func (h *Handler) ForgotPasswordHandler(c echo.Context) error {
	// bind request body
	payload := new(api_models.ForgotPasswordRequest)

	if err := h.BindBodyRequest(c, payload); err != nil {
		return api_models.SendBadRequestResponse(c, err.Error())
	}

	// validation
	validateErrors := h.ValidateBodyRequest(c, *payload)
	if validateErrors != nil {
		return api_models.SendFailedValidationResponse(c, "Request validation failed", validateErrors)
	}

	userService := services.NewUserService(h.DB)
	userRetrieved, err := userService.GetUserByEmail(payload.Email)
	if err != nil {
		h.Logger.Error(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api_models.SendNotFoundResponse(c, "No user found for that email")
		}
		return api_models.SendInternalServerErrorResponse(c, "Somehting went wrong while trying to fech user by email")
	}

	appTokenService := services.NewAppTokenService(h.DB)
	generatedAppToken, err := appTokenService.GenerateResetPasswordToken(*userRetrieved)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendInternalServerErrorResponse(c, "Something went wrong while trying to generate the reset password token")
	}
	encodedEmail := base64.RawURLEncoding.EncodeToString([]byte(userRetrieved.Email))
	fromtendURL, err := url.Parse(payload.FrontendURL)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendBadRequestResponse(c, "Invalid frontend URL")
	}
	queryUrl := url.Values{}
	queryUrl.Add("email", encodedEmail)
	queryUrl.Add("token", generatedAppToken.Token)
	fromtendURL.RawQuery = queryUrl.Encode()

	emailData := mailer.EmailData{
		Subject: "Did you forget your pass? No prob",
		Meta: struct {
			FrontendUrl string
			Token       string
		}{
			FrontendUrl: fromtendURL.String(),
			Token:       generatedAppToken.Token,
		},
	}
	err = h.Mailer.SendEmail(payload.Email, "password-reset.html", emailData)
	if err != nil {
		c.Logger().Error("Error sending password reset email:", err)
	}

	return api_models.SendSuccessResponse(c, "Email to reset the password has been sent.", nil)
}

func (h *Handler) ResetPasswordHandler(c echo.Context) error {
	// bind request body
	payload := new(api_models.ResetPasswordRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return api_models.SendBadRequestResponse(c, err.Error())
	}
	// validation
	validateErrors := h.ValidateBodyRequest(c, *payload)
	if validateErrors != nil {
		return api_models.SendFailedValidationResponse(c, "Request validation failed", validateErrors)
	}
	decodedEmail, err := base64.RawURLEncoding.DecodeString(payload.Meta)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendBadRequestResponse(c, "Invalid email encoding")
	}

	userService := services.NewUserService(h.DB)
	userRetrieved, err := userService.GetUserByEmail(string(decodedEmail))
	if err != nil {
		h.Logger.Error(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api_models.SendNotFoundResponse(c, "No user found for that email")
		}
		return api_models.SendInternalServerErrorResponse(c, "Somehting went wrong while trying to fech user by email")
	}

	appTokenService := services.NewAppTokenService(h.DB)
	validatedToken, err := appTokenService.ValidateResetPasswordToken(*userRetrieved, payload.Token)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendBadRequestResponse(c, err.Error())
	}
	err = userService.UpdateUserPassword(userRetrieved, payload.NewPassword)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendInternalServerErrorResponse(c, err.Error())
	}
	err = appTokenService.MarkTokenAsUsed(validatedToken)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendInternalServerErrorResponse(c, "Something went wrong while trying to mark the reset password token as used. "+err.Error())
	}
	return api_models.SendSuccessResponse(c, "Password reset successful", nil)
}
