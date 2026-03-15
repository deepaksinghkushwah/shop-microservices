package service

import (
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/repository"
)

func CreateProduct(name, slug, description string, price float64, categoryID uint) (*model.Product, error) {

	product := model.Product{
		Name:        name,
		Slug:        slug,
		Description: description,
		Price:       price,
		CategoryID:  categoryID,
	}

	err := repository.CreateProduct(&product)

	return &product, err
}

func ListProducts(page int, limit int, categoryID uint) ([]model.Product, int64, error) {

	return repository.GetProducts(page, limit, categoryID)
}

func GetProduct(slug string) (*model.Product, error) {

	return repository.GetProductBySlug(slug)
}

func GetAllProductImages() ([]model.ProductImage, error) {

	return repository.GetAllProductImages()
}
