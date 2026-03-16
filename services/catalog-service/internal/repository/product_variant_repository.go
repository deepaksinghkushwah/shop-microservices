package repository

import (
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/pkg/database"
)

func CreateVariant(variant *model.ProductVariant, valueIDs []uint) error {

	tx := database.DB.Begin()

	if err := tx.Create(variant).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, id := range valueIDs {

		mapping := model.VariantAttributeValue{
			ProductVariantID: variant.ID,
			AttributeValueID: id,
		}

		if err := tx.Create(&mapping).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func GetVariantsByProduct(productID uint) ([]model.ProductVariant, error) {

	var variants []model.ProductVariant

	err := database.DB.
		Where("product_id = ?", productID).
		Find(&variants).Error

	return variants, err
}
