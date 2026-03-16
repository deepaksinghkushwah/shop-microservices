package model

import "gorm.io/gorm"

type Attribute struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`
	Slug string `gorm:"uniqueIndex"`

	Values []AttributeValue `gorm:"foreignKey:AttributeID"`
}
