package handlers

import (
	"github.com/labstack/echo/v4"
	"go.budget-backend/cmd/api/api_models"
	"go.budget-backend/cmd/api/services"
	"go.budget-backend/common"
	"go.budget-backend/internal/models"
)

func (h *Handler) UpdateUserPassword(c echo.Context) error {
	userRetrieved, ok := c.Get("user").(models.UserModel)
	if !ok {
		h.Logger.Error("User not found in context")
		return api_models.SendUnauthorizedResponse(c, "User not authenticated")
	}

	payload := new(api_models.ChangePasswordRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return api_models.SendBadRequestResponse(c, err.Error())
	}

	if !common.CheckPasswordHash(payload.CurrentPassword, userRetrieved.Password) {
		return api_models.SendBadRequestResponse(c, "Invalid password.")
	}

	validationErr := api_models.ValidateBodyRequest(c, payload)
	if validationErr != nil {
		c.Logger().Error(validationErr)
		return api_models.SendFailedValidationResponse(c, "Validation failed", validationErr)
	}
	userService := services.NewUserService(h.DB)
	err := userService.UpdateUserPassword(&userRetrieved, payload.NewPassword)
	if err != nil {
		c.Logger().Error(err)
		return api_models.SendBadRequestResponse(c, "Password update failed")
	}

	return api_models.SendSuccessResponse(c, "Password updated successfully", nil)
}
