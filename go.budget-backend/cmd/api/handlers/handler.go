package handlers

import (
	"errors"

	"github.com/labstack/echo/v4"
	"go.budget-backend/internal/mailer"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Logger echo.Logger
	Mailer mailer.Mailer
}

func (h *Handler) BindBodyRequest(c echo.Context, payload interface{}) error {
	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		h.Logger.Error(err)
		return errors.New("Failed binding. " + err.Error())
	}
	return nil
}
