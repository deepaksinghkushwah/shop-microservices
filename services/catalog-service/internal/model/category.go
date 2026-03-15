package model

type Category struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Slug     string `gorm:"uniqueIndex"`
	ParentID *uint
}
