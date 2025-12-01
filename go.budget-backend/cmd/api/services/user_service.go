package services

import (
	"errors"
	"fmt"

	"go.budget-backend/cmd/api/api_models"
	"go.budget-backend/common"
	"go.budget-backend/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
	// Add any dependencies like database connections here
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (userService UserService) CreateUser(userRequest api_models.RegisterUserRequest) (*models.UserModel, error) {
	hashedPassword, err := common.HashPassword(userRequest.Password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return nil, errors.New("user creation failed")
	}

	user := models.UserModel{
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Email:     userRequest.Email,
		Password:  hashedPassword,
		Gender:    userRequest.Gender,
	}
	fmt.Println(user)

	result := userService.db.Create(&user)

	if result.Error != nil {
		fmt.Println("Error creating user:", result.Error)
		return nil, errors.New("user creation failed")
	}

	return &user, nil
}

func (userService UserService) GetUserByEmail(email string) (*models.UserModel, error) {
	var user models.UserModel
	result := userService.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (userService UserService) GetUserByID(email string) (*models.UserModel, error) {
	var user models.UserModel
	result := userService.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (userService UserService) UpdateUserPassword(user *models.UserModel, newPassword string) error {
	hashedPassword, err := common.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}
	result := userService.db.Model(user).Update("password", hashedPassword)
	if result.Error != nil {
		return errors.New("failed to update password")
	}
	return nil
}
