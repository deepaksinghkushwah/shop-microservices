package model

import "gorm.io/gorm"

type ProductImage struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	ProductID uint
	URL       string
	IsPrimary bool
}

type Product struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Slug        string `gorm:"uniqueIndex"`
	Description string
	Price       float64
	CategoryID  uint

	Images   []ProductImage   `gorm:"foreignKey:ProductID"`
	Variants []ProductVariant `gorm:"foreignKey:ProductID"`
}
