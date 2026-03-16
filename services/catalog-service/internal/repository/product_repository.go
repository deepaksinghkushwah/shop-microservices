package repository

import (
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/pkg/database"
)

func CreateProduct(product *model.Product) error {

	return database.DB.Create(product).Error
}

func GetProducts(page int, limit int, categoryID uint) ([]model.Product, int64, error) {

	var products []model.Product
	var total int64
	offset := (page - 1) * limit

	query := database.DB.Model(&model.Product{})
	if categoryID != 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	// Count total
	query.Count(&total)

	err := query.
		Preload("Images").
		Preload("Variants").
		Preload("Variants.Attributes.AttributeValue").
		Preload("Variants.Attributes.AttributeValue.Attribute").
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	return products, total, err
}

func GetProductBySlug(slug string) (*model.Product, error) {

	var product model.Product

	err := database.DB.
		Preload("Images").
		Preload("Variants").
		Preload("Variants.Attributes.AttributeValue").
		Preload("Variants.Attributes.AttributeValue.Attribute").
		Where("slug = ?", slug).
		First(&product).Error

	return &product, err
}
