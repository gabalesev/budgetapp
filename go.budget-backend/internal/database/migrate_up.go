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
	err = db.AutoMigrate(&models.UserModel{})
	if err != nil {
		panic(err)
	}
	log.Println("Migration completed.")
}
