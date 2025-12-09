package main

import (
	"fmt"

	"go.budget-backend/cmd/api/api_models"
	"go.budget-backend/cmd/api/services"
	"go.budget-backend/common"
)

func main() {
	db, err := common.NewMySQL()
	if err != nil {
		panic(err)
	}
	categoryService := services.CategoryService{
		DB: db,
	}
	categories := []string{
		"Food",
		"Transportation",
		"Utilities",
		"Entertainment",
		"Healthcare",
		"Education",
		"Shopping",
		"Travel",
		"Personal Care",
		"Miscellaneous",
	}

	for _, category := range categories {
		_, err := categoryService.Create(&api_models.CategoryRequest{
			Name:     category,
			IsCustom: false,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("Catyegory created: ", category)
	}

}
