package services

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"go.budget-backend/internal/models"
	"gorm.io/gorm"
)

type AppTokenService struct {
	db *gorm.DB
	// Add any dependencies like database connections here
}

func NewAppTokenService(db *gorm.DB) *AppTokenService {
	return &AppTokenService{db: db}
}

func (appTokenService AppTokenService) getToken() int {
	rand.Seed(time.Now().UnixNano())
	min := 10000
	max := 99999
	return rand.Intn(max-min+1) + min
}

func (appTokenService AppTokenService) GenerateResetPasswordToken(user models.UserModel) (*models.AppTokenModel, error) {
	token := appTokenService.getToken()
	appToken := models.AppTokenModel{
		TargetID:  user.ID,
		Token:     strconv.Itoa(token),
		Type:      "reset_password",
		ExpiredAt: time.Now().Add(15 * time.Minute),
		Used:      false,
	}
	result := appTokenService.db.Create(&appToken)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println("Generated token:", appToken.Token)
	return &appToken, nil
}

func (appTokenService AppTokenService) ValidateResetPasswordToken(user models.UserModel, token string) (*models.AppTokenModel, error) {
	var retrievedToken models.AppTokenModel
	result := appTokenService.db.Where(&models.AppTokenModel{
		TargetID: user.ID,
		Token:    token,
		Type:     "reset_password",
	}).First(&retrievedToken)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid password reset token")
		}
		return nil, result.Error
	}
	if retrievedToken.Used {
		return nil, errors.New("token already used")
	}
	if retrievedToken.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("token has expired. Reinitiate password reset process")
	}

	return &retrievedToken, nil
}

func (appTokenService AppTokenService) MarkTokenAsUsed(token *models.AppTokenModel) error {
	token.Used = true
	result := appTokenService.db.Save(token)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
