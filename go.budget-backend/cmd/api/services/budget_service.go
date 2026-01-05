package services

import (
	"errors"
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

func (c BudgetService) ListBudgets(query *gorm.DB, pagination *common.Pagination, budgets []*models.BudgetModel) (*common.Pagination, error) {

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

	budgetCount, err := b.GetBudgetCount(model.UserID, model.Month, model.Year, model.Slug)
	if err == nil && budgetCount != 0 {
		return nil, errors.New("category with this name already exists or existed")
	}

	result := b.DB.Create(model)
	if result.Error != nil {
		return nil, errors.New("failed to create budget: " + result.Error.Error())
	}
	return model, nil
}

func (b BudgetService) GetByID(userID uint, budgetID uint) (*models.BudgetModel, error) {
	var budgetRetrieved models.BudgetModel
	result := b.DB.
		Where("user_id = ? AND id = ?", userID, budgetID).
		First(&budgetRetrieved)
	if result.Error != nil {
		return nil, result.Error
	}
	return &budgetRetrieved, nil
}

func (b BudgetService) GetBudgetCount(userID uint, month uint, year uint16, slug string) (int64, error) {
	var budgetCount int64
	result := b.DB.Model(models.BudgetModel{}).
		Where("user_id = ? AND month = ? AND year = ? AND slug = ?", userID, month, year, slug).
		Count(&budgetCount)
	if result.Error != nil {
		return -1, result.Error
	}
	return budgetCount, nil
}

func (b BudgetService) Update(budget *models.BudgetModel, payload *api_models.UpdateBudgetRequest, userID uint) (*models.BudgetModel, error) {
	if payload.Date != "" {
		timeParsed, err := time.Parse(time.DateOnly, payload.Date)
		if err != nil {
			return nil, err
		}
		budget.Date = timeParsed
	}

	if payload.Amount != 0 {
		budget.Amount = payload.Amount
	}

	if payload.Title != "" {
		budget.Title = payload.Title
		slug := strings.ToLower(payload.Title)
		budget.Slug = strings.Replace(slug, " ", "_", -1)
	}

	budgetCount, _ := b.GetBudgetCount(budget.UserID, uint(budget.Date.Month()), uint16(budget.Date.Year()), budget.Slug)
	if budgetCount > 0 {
		return nil, errors.New("Budget with selected month, year and title already exists.")
	}

	if payload.Description != "" {
		budget.Description = &payload.Description
	}

	result := b.DB.Save(budget)
	if result.Error != nil {
		return nil, errors.New(result.Error.Error())
	}

	return budget, nil
}
