package common

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMySQL() (*gorm.DB, error) {
	err := godotenv.Load("internal/.env")
	if err != nil {
		return nil, err
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", username, password, host, port, database)
	fmt.Printf("DB: %s", dns)
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger:         logger.Default,
		TranslateError: true,
	})

	if err != nil {
		return nil, err
	}
	log.Default().Println("Database connection successful")
	return db, nil
}
