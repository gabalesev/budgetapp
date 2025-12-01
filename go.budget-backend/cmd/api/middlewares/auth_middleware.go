package middlewares

import (
	"errors"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"go.budget-backend/cmd/api/api_models"
	"go.budget-backend/common"
	"go.budget-backend/internal/models"
	"gorm.io/gorm"
)

type AppMiddleware struct {
	Logger echo.Logger
	DB     *gorm.DB
}

func (appMiddleware *AppMiddleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
		c.Response().Header().Add("Vary", "Authorization")
		tokenWithPrefix := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(tokenWithPrefix, "Bearer ") {
			return api_models.SendUnauthorizedResponse(c, "Missing bearer token")
		}
		accessToken := strings.TrimPrefix(tokenWithPrefix, "Bearer ")

		claims, err := common.ParseJWT(accessToken)
		if err != nil {
			return api_models.SendUnauthorizedResponse(c, "Invalid or expired token")
		}

		if common.IsTokenExpired(*claims) {
			return api_models.SendUnauthorizedResponse(c, "Token has expired")
		}
		var user models.UserModel
		result := appMiddleware.DB.First(&user, claims.ID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return api_models.SendUnauthorizedResponse(c, "Invalid token")
		}

		if result.Error != nil {
			appMiddleware.Logger.Error("Error fetching user: ", result.Error)
			return api_models.SendUnauthorizedResponse(c, "User not found")
		}

		fmt.Println("Authorization header:", tokenWithPrefix)
		fmt.Println(user)

		c.Set("user", user)
		return next(c)
	}
}
