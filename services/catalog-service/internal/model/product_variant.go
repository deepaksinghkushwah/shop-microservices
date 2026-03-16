package model

import "gorm.io/gorm"

type ProductVariant struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	ProductID uint
	SKU       string `gorm:"uniqueIndex"`
	Price     float64
	Stock     int

	Product    Product
	Attributes []VariantAttributeValue `gorm:"foreignKey:ProductVariantID"`
}
