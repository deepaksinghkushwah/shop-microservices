package service

import (
	"errors"
	"strings"

	"github.com/deepaksinghkushwah/shop-microservices/pkg/auth"
	"github.com/deepaksinghkushwah/shop-microservices/services/auth-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/auth-service/internal/repository"
	"github.com/deepaksinghkushwah/shop-microservices/services/auth-service/pkg/hash"
)

func Register(name, email, password string) error {
	email = strings.ToLower(email)
	hashedPassword, err := hash.HashPassword(password)
	if err != nil {
		return err
	}
	user := model.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	return repository.CreateUser(&user)
}

func Login(email, password string) (string, error) {
	email = strings.ToLower(email)
	user, err := repository.FindByEmail(email)
	if err != nil {
		return "", err
	}

	if !hash.CheckPassword(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return auth.GenerateToken(user.ID)
}

func Profile(userID uint) (*model.User, error) {
	return repository.FindByID(userID)
}
