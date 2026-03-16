package repository

import (
	"errors"

	"github.com/deepaksinghkushwah/shop-microservices/services/order-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/order-service/pkg/database"
	"gorm.io/gorm"
)

func CreateOrderItem(orderItem *model.OrderItem) error {
	if err := database.DB.Create(orderItem).Error; err != nil {
		return err
	}
	return nil
}

func GetOrderItem(itemID uint) (*model.OrderItem, error) {
	var orderItem model.OrderItem
	if err := database.DB.First(&orderItem, itemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order item not found")
		}
		return nil, err
	}
	return &orderItem, nil
}

func GetOrderItemsByOrder(orderID uint) ([]model.OrderItem, error) {
	var orderItems []model.OrderItem
	if err := database.DB.Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
		return nil, err
	}
	return orderItems, nil
}

func UpdateOrderItem(itemID uint, quantity int, totalPrice float64) (*model.OrderItem, error) {
	orderItem, err := GetOrderItem(itemID)
	if err != nil {
		return nil, err
	}

	if err := database.DB.Model(orderItem).Updates(map[string]interface{}{
		"quantity":    quantity,
		"total_price": totalPrice,
	}).Error; err != nil {
		return nil, err
	}

	return orderItem, nil
}

func DeleteOrderItem(itemID uint) error {
	orderItem, err := GetOrderItem(itemID)
	if err != nil {
		return err
	}

	if err := database.DB.Delete(orderItem).Error; err != nil {
		return err
	}

	return nil
}
