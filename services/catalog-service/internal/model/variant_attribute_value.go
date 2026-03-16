package model

type VariantAttributeValue struct {
	ID uint `gorm:"primaryKey"`

	ProductVariantID uint
	AttributeValueID uint

	// swaggerignore: true
	ProductVariant ProductVariant `json:"-" swaggerignore:"true"`

	// swaggerignore: true
	AttributeValue AttributeValue `json:"-" swaggerignore:"true"`
}
