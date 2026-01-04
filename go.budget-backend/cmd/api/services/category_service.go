package services

import (
	"errors"
	"fmt"
	"strings"

	"go.budget-backend/cmd/api/api_models"
	"go.budget-backend/common"
	"go.budget-backend/internal/models"
	"gorm.io/gorm"
)

type CategoryService struct {
	DB *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{DB: db}
}

func (c CategoryService) GetAllCategories(pagination *common.Pagination, categories []*models.CategoryModel) (*common.Pagination, error) {

	//result := c.DB.Find(&categories)
	err := c.DB.Scopes(pagination.Paginate()).Find(&categories)
	if err.Error != nil {
		return nil, errors.New("failed to fetch categories")
	}
	pagination.Items = categories

	return pagination, nil
}

func (c CategoryService) GetMultipleCategories(categoryIDs []uint) ([]*models.CategoryModel, error) {
	var categories []*models.CategoryModel
	err := c.DB.Where("id IN ?", categoryIDs).Find(&categories)
	if err.Error != nil {
		return nil, errors.New("failed to fetch categories")
	}

	return categories, nil
}

func (c CategoryService) Create(data *api_models.CategoryRequest) (*models.CategoryModel, error) {
	slug := strings.ToLower(data.Name)
	slug = strings.Replace(slug, " ", "_", -1)
	categoryCreated := &models.CategoryModel{
		Name:     data.Name,
		Slug:     slug,
		IsCustom: data.IsCustom,
	}
	result := c.DB.Where(models.CategoryModel{Slug: slug, Name: data.Name}).First(categoryCreated)

	if result.Error != nil {
		return nil, errors.New("category with this name already exists or existed")
	}

	result = c.DB.Create(categoryCreated)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, errors.New("failed to create category: " + result.Error.Error())
	}
	return categoryCreated, nil
}

func (c CategoryService) DeleteByID(category *models.CategoryModel) error {
	result := c.DB.Delete(&category)
	if result.Error != nil {
		return errors.New("failed to delete category")
	}
	return nil
}

func (c CategoryService) GetByID(categoryID uint) (*models.CategoryModel, error) {
	var category *models.CategoryModel
	result := c.DB.First(&category, categoryID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, errors.New("failed to fetch category")
	}
	return category, nil
}
