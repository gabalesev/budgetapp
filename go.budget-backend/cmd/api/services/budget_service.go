package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.budget-backend/cmd/api/api_models"
	"go.budget-backend/common"
	"go.budget-backend/internal/models"
	"gorm.io/gorm"
)

type BudgetService struct {
	DB *gorm.DB
}

func NewBudgetService(db *gorm.DB) *BudgetService {
	return &BudgetService{DB: db}
}

func (c BudgetService) GetAllBudgets(query *gorm.DB, pagination *common.Pagination, budgets []*models.BudgetModel) (*common.Pagination, error) {

	//result := c.DB.Find(&categories)
	err := query.Scopes(pagination.Paginate()).Find(&budgets)
	if err.Error != nil {
		return nil, errors.New("failed to fetch categories")
	}
	pagination.Items = budgets

	return pagination, nil
}

func (b BudgetService) Create(payload *api_models.CreateBudgetRequest, UserID uint) (*models.BudgetModel, error) {
	slug := strings.ToLower(payload.Title)
	slug = strings.Replace(slug, " ", "_", -1)

	model := &models.BudgetModel{
		Title:       payload.Title,
		Amount:      payload.Amount,
		UserID:      UserID,
		Slug:        slug,
		Description: &payload.Description,
	}

	if payload.Date == "" {
		currentDate := time.Now()
		model.Date = currentDate
	}

	model.Month = uint(model.Date.Month())
	model.Year = uint16(model.Date.Year())

	budgetRetrieved, err := b.budgetExists(model.UserID, model.Month, model.Year, model.Slug)
	if err == nil {

		return budgetRetrieved, errors.New("category with this name already exists or existed")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	result := b.DB.Create(model)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, errors.New("failed to create budget: " + result.Error.Error())
	}
	return model, nil
}

func (b BudgetService) budgetExists(userID uint, month uint, year uint16, slug string) (*models.BudgetModel, error) {
	var budgetRetrieved models.BudgetModel
	result := b.DB.
		Where("user_id = ? AND month = ? AND year = ? AND slug = ?", userID, month, year, slug).
		First(&budgetRetrieved)
	if result.Error != nil {
		return nil, result.Error
	}
	return &budgetRetrieved, nil
}
