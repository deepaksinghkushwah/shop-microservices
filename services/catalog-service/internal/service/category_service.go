package service

import (
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/repository"
)

func CreateCategory(name string, slug string, parentID *uint) (*model.Category, error) {

	category := model.Category{
		Name:     name,
		Slug:     slug,
		ParentID: parentID,
	}

	err := repository.CreateCategory(&category)

	return &category, err
}

func ListCategories() ([]model.Category, error) {

	return repository.GetCategories()
}
