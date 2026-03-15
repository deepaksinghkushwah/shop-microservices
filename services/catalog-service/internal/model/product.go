package model

type ProductImage struct {
	ID        uint `gorm:"primaryKey"`
	ProductID uint
	URL       string
	IsPrimary bool
}

type Product struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Slug        string `gorm:"uniqueIndex"`
	Description string
	Price       float64
	CategoryID  uint

	Images []ProductImage `gorm:"foreignKey:ProductID"`
}
