package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model

	UserID     *uint       `json:"user_id"` // nil for guest orders
	Status     string      `json:"status" gorm:"not null;default:'pending'"`
	Note       string      `json:"note" gorm:"type:text"` // special requests
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
}
