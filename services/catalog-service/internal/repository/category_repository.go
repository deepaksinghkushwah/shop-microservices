package repository

import (
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/pkg/database"
)

func CreateCategory(category *model.Category) error {

	return database.DB.Create(category).Error
}

func GetCategories() ([]model.Category, error) {

	var categories []model.Category

	err := database.DB.Find(&categories).Error

	return categories, err
}
