package repository

import (
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/pkg/database"
)

func AssignAttributeToProduct(productID uint, attributeID uint) error {

	productAttribute := model.ProductAttribute{
		ProductID:   productID,
		AttributeID: attributeID,
	}

	return database.DB.Create(&productAttribute).Error
}
