package service

import (
	"github.com/deepaksinghkushwah/shop-microservices/services/order-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/order-service/internal/repository"
)

func CreateOrder(userID *uint, note string, items []model.OrderItem) (*model.Order, error) {
	order := &model.Order{
		UserID:     userID,
		Status:     "pending",
		Note:       note,
		OrderItems: items,
	}

	if err := repository.CreateOrder(order); err != nil {
		return nil, err
	}

	return order, nil
}

func GetOrdersByUser(userID uint) ([]model.Order, error) {
	return repository.GetOrdersByUser(userID)
}

func UpdateOrderStatus(orderID uint, status string) error {
	return repository.UpdateOrderStatus(orderID, status)
}
