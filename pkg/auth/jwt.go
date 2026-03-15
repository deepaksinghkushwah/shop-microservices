package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint) (string, error) {

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret" // fallback for testing
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string) (uint, error) {

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret" // fallback for testing
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)

	userID := uint(claims["user_id"].(float64))

	return userID, nil
}
