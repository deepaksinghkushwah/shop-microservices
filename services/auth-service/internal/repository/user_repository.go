package repository

import (
	"github.com/deepaksinghkushwah/shop-microservices/services/auth-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/auth-service/pkg/database"
)

func CreateUser(user *model.User) error {
	return database.DB.Create(user).Error
}

func FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func FindByID(id uint) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, id).Error
	return &user, err
}
