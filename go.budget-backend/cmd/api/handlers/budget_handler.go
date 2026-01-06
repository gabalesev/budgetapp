package handlers

import (
	"errors"

	"github.com/labstack/echo/v4"
	"go.budget-backend/cmd/api/api_models"
	"go.budget-backend/cmd/api/services"
	"go.budget-backend/common"
	"go.budget-backend/internal/models"
	"gorm.io/gorm"
)

func (h *Handler) GetAllBudgetsHandler(c echo.Context) error {
	budgetService := services.NewBudgetService(h.DB)
	var budgets []*models.BudgetModel
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		c.Logger().Error("User not found in context")
		return api_models.SendInternalServerErrorResponse(c, "Something went wrong while fetching user.")
	}
	query := h.DB.Preload("Categories").Scopes(common.WhereUserIDScope(user.ID))

	paginator := common.NewPagination(budgets, c.Request(), query)

	paginatedBudgets, err := budgetService.ListBudgets(query, paginator, budgets)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendInternalServerErrorResponse(c, err.Error())
	}
	return api_models.SendSuccessResponse(c, "Budgets retrieved successfully", paginatedBudgets)
}

func (h *Handler) CreateBudgetHandler(c echo.Context) error {

	payload := new(api_models.CreateBudgetRequest)

	err := h.BindBodyRequest(c, payload)
	if err != nil {
		return api_models.SendBadRequestResponse(c, "Invalid request payload")
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		h.Logger.Error(err)
		return api_models.SendFailedValidationResponse(c, "Invalid request payload", validationErrors)
	}
	categoryService := services.CategoryService{DB: h.DB}
	categories, err := categoryService.GetMultipleCategories(payload.Categories)
	if err != nil {
		return api_models.SendInternalServerErrorResponse(c, "Failed to fetch categories")
	}

	user, _ := c.Get("user").(models.UserModel)
	createdBudget := &models.BudgetModel{}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		budgetService := services.BudgetService{DB: tx}

		createdBudget, err = budgetService.Create(payload, user.ID)
		if err != nil {
			return errors.New("Budget could  not be created")
		}

		err = tx.Model(createdBudget).Association("Categories").Replace(categories)
		if err != nil {
			return errors.New("Failed to associate categories with budget")
		}
		createdBudget.Categories = categories

		return nil
	})

	if err != nil {
		return api_models.SendInternalServerErrorResponse(c, err.Error())
	}

	return api_models.SendSuccessResponse(c, "Budget created successfully", createdBudget)
}

func (h *Handler) UpdateBudgetHandler(c echo.Context) error {
	var budgetId api_models.IDParamRequest
	err := (&echo.DefaultBinder{}).BindPathParams(c, &budgetId)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendBadRequestResponse(c, err.Error())
	}

	budgetService := services.NewBudgetService(h.DB)
	categoryService := services.NewCategoryService(h.DB)

	// bind request body
	payload := new(api_models.UpdateBudgetRequest)
	er := h.BindBodyRequest(c, payload)
	if er != nil {
		return api_models.SendBadRequestResponse(c, er.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		h.Logger.Error(er)
		return api_models.SendFailedValidationResponse(c, "Invalid request payload", validationErrors)
	}

	user, _ := c.Get("user").(models.UserModel)

	budget, err := budgetService.GetByID(user.ID, budgetId.ID)
	if err != nil || user.ID != budget.UserID {
		if err != nil {
			h.Logger.Error(err)
		}
		return api_models.SendNotFoundResponse(c, "Budget not found")
	}

	updatedBudget := &models.BudgetModel{}
	var categories []*models.CategoryModel
	if payload.Categories != nil {
		categories, err = categoryService.GetMultipleCategories(payload.Categories)
		if err != nil {
			h.Logger.Error(err)
			return api_models.SendInternalServerErrorResponse(c, "Failed to fetch categories")
		}
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		budgetService := services.BudgetService{DB: tx}

		updatedBudget, err = budgetService.Update(budget, payload, budgetId.ID)

		if err != nil {
			h.Logger.Error(err)
			return errors.New(err.Error())
		}
		if len(categories) > 0 {
			err = tx.Model(&updatedBudget).Association("Categories").Replace(categories)
			if err != nil {
				h.Logger.Error(err)
				return errors.New("Failed to associate categories with budget")
			}
			updatedBudget.Categories = categories
		}

		return nil
	})

	if err != nil {
		if err.Error() == "Budget with selected month, year and title already exists." {
			return api_models.SendBadRequestResponse(c, err.Error())
		}
		return api_models.SendInternalServerErrorResponse(c, err.Error())
	}

	return api_models.SendSuccessResponse(c, "Budget updated successfully", updatedBudget)
}

func (h *Handler) DeleteBudgetHandler(c echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		h.Logger.Error("User not found in context")
		return api_models.SendInternalServerErrorResponse(c, "Something went wrong while fetching user.")
	}

	var budgetId api_models.IDParamRequest
	err := (&echo.DefaultBinder{}).BindPathParams(c, &budgetId)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendBadRequestResponse(c, err.Error())
	}

	budgetService := services.NewBudgetService(h.DB)

	budget, err := budgetService.GetByID(user.ID, budgetId.ID)
	if err != nil || user.ID != budget.UserID {
		if err != nil {
			h.Logger.Error(err)
		}
		return api_models.SendNotFoundResponse(c, "Budget not found")
	}

	err = h.DB.Model(&budget).Association("Categories").Clear()
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendInternalServerErrorResponse(c, "Failed to clear budget categories")
	}

	query := h.DB.Scopes(common.WhereUserIDScope(user.ID))
	err = query.Delete(&models.BudgetModel{}, budget.ID).Error
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendInternalServerErrorResponse(c, "Failed to delete budget")
	}
	return api_models.SendSuccessResponse(c, "Budget deleted successfully", nil)
}
