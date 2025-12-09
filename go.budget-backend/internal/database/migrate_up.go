package main

import (
	"log"

	"go.budget-backend/common"
	"go.budget-backend/internal/models"
)

func main() {
	db, err := common.NewMySQL()
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(
		&models.UserModel{},
		&models.AppTokenModel{},
		&models.CategoryModel{},
		&models.BudgetModel{})
	if err != nil {
		panic(err)
	}
	log.Println("Migration completed.")
}
