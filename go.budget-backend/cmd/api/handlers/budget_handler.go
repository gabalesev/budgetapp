package handlers

import (
	"github.com/labstack/echo/v4"
	"go.budget-backend/cmd/api/api_models"
	"go.budget-backend/cmd/api/services"
	"go.budget-backend/common"
	"go.budget-backend/internal/models"
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

	budgetService := services.BudgetService{
		DB: h.DB,
	}
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		h.Logger.Error("User not found in context")
		return api_models.SendInternalServerErrorResponse(c, "Something went wrong while fetching user.")
	}
	createdBudget, err := budgetService.Create(payload, user.ID)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendInternalServerErrorResponse(c, "Failed to create budget")
	}

	categoryService := services.CategoryService{
		DB: h.DB,
	}

	categories, err := categoryService.GetMultipleCategories(payload.Categories)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendInternalServerErrorResponse(c, "Failed to fetch categories")
	}

	err = h.DB.Model(&createdBudget).Association("Categories").Replace(categories)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendInternalServerErrorResponse(c, "Failed to associate categories with budget")
	}

	createdBudget.Categories = categories

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

	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		h.Logger.Error("User not found in context")
		return api_models.SendInternalServerErrorResponse(c, "Something went wrong while fetching user.")
	}

	budget, err := budgetService.GetByID(user.ID, budgetId.ID)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendInternalServerErrorResponse(c, "Failed to fetch budget")
	}

	if user.ID != budget.UserID {
		h.Logger.Error("Unauthorized access to budget")
		return api_models.SendErrorResponse(c, "Unauthorized access to budget", 403)
	}

	createdBudget, err := budgetService.Update(budget, payload, budgetId.ID)
	if err != nil {
		h.Logger.Error(err)
		if err.Error() == "Budget with selected month, year and title already exists." {
			return api_models.SendBadRequestResponse(c, err.Error())
		}
		return api_models.SendInternalServerErrorResponse(c, "Failed to update budget")
	}

	categories, err := categoryService.GetMultipleCategories(payload.Categories)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendInternalServerErrorResponse(c, "Failed to fetch categories")
	}

	err = h.DB.Model(&createdBudget).Association("Categories").Replace(categories)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendInternalServerErrorResponse(c, "Failed to associate categories with budget")
	}

	createdBudget.Categories = categories

	return api_models.SendSuccessResponse(c, "Budget updated successfully", createdBudget)
}
