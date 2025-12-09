package handlers

import (
	"github.com/labstack/echo/v4"
	"go.budget-backend/cmd/api/api_models"
	"go.budget-backend/cmd/api/services"
	"go.budget-backend/common"
	"go.budget-backend/internal/models"
)

func (h *Handler) GetAllCategoriesHandler(c echo.Context) error {
	categoryService := services.NewCategoryService(h.DB)
	var categories []*models.CategoryModel
	paginator := common.NewPagination(categories, c.Request(), h.DB)

	paginatedCategories, err := categoryService.GetAllCategories(paginator, categories)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendInternalServerErrorResponse(c, err.Error())
	}
	return api_models.SendSuccessResponse(c, "Categories retrieved successfully", paginatedCategories)
}

func (h *Handler) CreateCategoryHandler(c echo.Context) error {

	payload := new(api_models.CategoryRequest)

	err := h.BindBodyRequest(c, payload)
	if err != nil {
		return api_models.SendBadRequestResponse(c, "Invalid request payload")
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		h.Logger.Error(err)
		return api_models.SendFailedValidationResponse(c, "Invalid request payload", validationErrors)
	}

	categoryService := services.CategoryService{
		DB: h.DB,
	}

	createdCategory, err := categoryService.Create(payload)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendInternalServerErrorResponse(c, err.Error())
	}
	return api_models.SendSuccessResponse(c, "Category created successfully", createdCategory)
}

func (h *Handler) GetCategoryByIDHandler(c echo.Context) error {

	var categoryID api_models.IDParamRequest
	err := (&echo.DefaultBinder{}).BindPathParams(c, &categoryID)
	if err != nil {
		return api_models.SendBadRequestResponse(c, "Invalid request parameters")
	}

	categoryService := services.CategoryService{
		DB: h.DB,
	}
	category, err := categoryService.GetByID(categoryID.ID)

	if err != nil {
		h.Logger.Error(err)
		return api_models.SendNotFoundResponse(c, err.Error())
	}

	return api_models.SendSuccessResponse(c, "Category retrieved successfully", category)

}

func (h *Handler) DeleteCategoryByIDHandler(c echo.Context) error {

	// TODO authorize by roles ?
	// _, ok := c.Get("user").(models.UserModel)
	// if !ok {
	// 	h.Logger.Error("User not found in context")
	// 	return api_models.SendUnauthorizedResponse(c, "Unauthorized")
	// }

	var categoryID api_models.IDParamRequest
	err := (&echo.DefaultBinder{}).BindPathParams(c, &categoryID)
	if err != nil {
		return api_models.SendBadRequestResponse(c, "Invalid request parameters")
	}

	categoryService := services.CategoryService{
		DB: h.DB,
	}

	category, err := categoryService.GetByID(categoryID.ID)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendNotFoundResponse(c, err.Error())
	}

	err = categoryService.DeleteByID(category)
	if err != nil {
		h.Logger.Error(err)
		return api_models.SendInternalServerErrorResponse(c, err.Error())
	}

	return api_models.SendSuccessResponse(c, "Category deleted", category)

}

func (h *Handler) UpdateCategoryHandler(c echo.Context) error {
	_, ok := c.Get("user").(models.UserModel)
	if !ok {
		h.Logger.Error("User not found in context")
		return api_models.SendUnauthorizedResponse(c, "Unauthorized")
	}
	return nil
}
