package model

import "gorm.io/gorm"

type ProductAttribute struct {
	gorm.Model

	ProductID   uint
	AttributeID uint

	Product   Product
	Attribute Attribute
}
