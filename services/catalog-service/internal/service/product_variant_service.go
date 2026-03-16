package service

import (
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/repository"
)

func CreateVariant(productID uint, sku string, price float64, stock int, attributeValueIDs []uint) (*model.ProductVariant, error) {

	variant := model.ProductVariant{
		ProductID: productID,
		SKU:       sku,
		Price:     price,
		Stock:     stock,
	}

	err := repository.CreateVariant(&variant, attributeValueIDs)

	return &variant, err
}
