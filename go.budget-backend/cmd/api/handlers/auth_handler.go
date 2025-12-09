package handlers

import (
	"errors"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"go.budget-backend/cmd/api/api_models"
	"go.budget-backend/cmd/api/services"
	"go.budget-backend/common"
	"go.budget-backend/internal/mailer"
	"go.budget-backend/internal/models"
	"gorm.io/gorm"
)

func (h *Handler) RegisterHandler(c echo.Context) error {
	// bind request body to struct
	payload := new(api_models.RegisterUserRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return api_models.SendBadRequestResponse(c, err.Error())
	}

	fmt.Println(*payload)

	validationErr := h.ValidateBodyRequest(c, payload)

	if validationErr != nil {
		c.Logger().Error(validationErr)
		return api_models.SendFailedValidationResponse(c, "Validation failed", validationErr)
	}

	userService := services.NewUserService(h.DB)

	_, err := userService.GetUserByEmail(*payload.Email)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.Logger().Error(err)
		return api_models.SendBadRequestResponse(c, "User already exists with this email")
	}
	user, err := userService.CreateUser(*payload)

	if err != nil {
		c.Logger().Error(err)
		return api_models.SendBadRequestResponse(c, "User creation failed")
	}
	emailData := mailer.EmailData{
		Subject: "Welcome to the " + os.Getenv("APP_NAME"),
		Meta: struct {
			FirstName string
			LoginLink string
		}{
			FirstName: *user.FirstName,
			LoginLink: "http://localhost:3000/login",
		},
	}
	err = h.Mailer.SendEmail(*payload.Email, "welcome.html", emailData)
	if err != nil {
		c.Logger().Error("Error sending welcome email:", err)
	}
	return api_models.SendSuccessResponse(c, "User registration successfull", user)
}

func (h *Handler) GetUserHandler(c echo.Context) error {
	email := c.Param("email")
	userService := services.NewUserService(h.DB)
	user, err := userService.GetUserByEmail(email)
	if err != nil {
		c.Logger().Error(err)
		return api_models.SendBadRequestResponse(c, "User not found")
	}
	return api_models.SendSuccessResponse(c, "User fetched successfully", user)
}

func (h *Handler) LoginHandler(c echo.Context) error {
	// bind request body to struct
	payload := new(api_models.LoginRequest)
	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		c.Logger().Error(err)
		return api_models.SendBadRequestResponse(c, err.Error())
	}

	validationErr := h.ValidateBodyRequest(c, payload)
	if validationErr != nil {
		c.Logger().Error(validationErr)
		return api_models.SendFailedValidationResponse(c, "Validation failed", validationErr)
	}

	userService := services.NewUserService(h.DB)
	userRetrieved, err := userService.GetUserByEmail(*payload.Email)
	if err != nil {
		c.Logger().Error(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api_models.SendBadRequestResponse(c, "Invalid email or password.")
		}
		return api_models.SendInternalServerErrorResponse(c, "Something went wrong while authenticating user.")
	}

	if !common.CheckPasswordHash(payload.Password, userRetrieved.Password) {
		return api_models.SendBadRequestResponse(c, "Invalid email or password.")
	}

	accessToken, refreshToken, err := common.GenerateJWT(*userRetrieved)
	if err != nil {
		c.Logger().Error(err)
		return api_models.SendInternalServerErrorResponse(c, "Something went wrong while generating tokens.")
	}

	return api_models.SendSuccessResponse(c, "User logged in", map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          userRetrieved,
	})
}

func (h Handler) GetAuthenticatedUserHandler(c echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		c.Logger().Error("User not found in context")
		return api_models.SendInternalServerErrorResponse(c, "Something went wrong while fetching user.")
	}
	return api_models.SendSuccessResponse(c, "User fetched successfully", user)
}
