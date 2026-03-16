package model

import "gorm.io/gorm"

type AttributeValue struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	AttributeID uint
	Value       string

	Attribute Attribute `json:"-"`
}
