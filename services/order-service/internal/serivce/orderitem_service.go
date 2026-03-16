package service

import (
	"github.com/deepaksinghkushwah/shop-microservices/services/order-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/order-service/internal/repository"
)

func CreateOrderItem(orderItem *model.OrderItem) (*model.OrderItem, error) {
	if err := repository.CreateOrderItem(orderItem); err != nil {
		return nil, err
	}
	return orderItem, nil
}

func GetOrderItem(itemID uint) (*model.OrderItem, error) {
	return repository.GetOrderItem(itemID)
}

func GetOrderItemsByOrder(orderID uint) ([]model.OrderItem, error) {
	return repository.GetOrderItemsByOrder(orderID)
}

func UpdateOrderItem(itemID uint, quantity int, totalPrice float64) (*model.OrderItem, error) {
	return repository.UpdateOrderItem(itemID, quantity, totalPrice)
}

func DeleteOrderItem(itemID uint) error {
	return repository.DeleteOrderItem(itemID)
}
