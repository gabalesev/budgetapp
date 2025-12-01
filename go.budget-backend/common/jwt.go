package common

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.budget-backend/internal/models"
)

type JWTClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.UserModel) (*string, *string, error) {
	userClaims := JWTClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, nil, err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 6)),
		},
	})
	signedRefreshToken, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, nil, err
	}
	return &signedToken, &signedRefreshToken, nil
}

func ParseJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func IsTokenExpired(claims JWTClaims) bool {
	currentTime := jwt.NewNumericDate(time.Now())
	return claims.ExpiresAt.Time.Before(currentTime.Time)
}
