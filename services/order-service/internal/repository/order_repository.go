package repository

import (
	"github.com/deepaksinghkushwah/shop-microservices/services/order-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/order-service/pkg/database"
)

func CreateOrder(order *model.Order) error {
	return database.DB.Create(order).Error
}

func GetOrdersByUser(userID uint) ([]model.Order, error) {
	var orders []model.Order
	err := database.DB.Preload("OrderItems").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func UpdateOrderStatus(orderID uint, status string) error {
	return database.DB.Model(&model.Order{}).Where("id = ?", orderID).Update("status", status).Error
}
