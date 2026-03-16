package model

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model

	OrderID      uint    `json:"order_id"`
	ProductID    uint    `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductSlug  string  `json:"product_slug"`
	VariantID    uint    `json:"variant_id"`
	VariantSKU   string  `json:"variant_sku"`
	VariantPrice float64 `json:"variant_price"`
	Quantity     int     `json:"quantity"`
	TotalPrice   float64 `json:"total_price"`
}
