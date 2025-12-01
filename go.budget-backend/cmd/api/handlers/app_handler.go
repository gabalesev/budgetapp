package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type helthCheck struct {
	Health bool `json:"health"`
}

func (h *Handler) HealthCheck(c echo.Context) error {
	healthCheckStruct := helthCheck{
		Health: true,
	}

	return c.JSON(http.StatusOK, healthCheckStruct)
}
