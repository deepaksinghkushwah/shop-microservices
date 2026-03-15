package repository

import (
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/pkg/database"
)

func CreateProductImage(image *model.ProductImage) error {

	return database.DB.Create(image).Error
}

func GetProductImages(productID uint) ([]model.ProductImage, error) {

	var images []model.ProductImage

	err := database.DB.
		Where("product_id = ?", productID).
		Find(&images).Error

	return images, err
}

func GetAllProductImages() ([]model.ProductImage, error) {

	var images []model.ProductImage

	err := database.DB.Find(&images).Error

	return images, err
}
