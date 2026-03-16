package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Slug     string `gorm:"uniqueIndex"`
	ParentID *uint
}
