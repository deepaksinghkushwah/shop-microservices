package handler

import (
	"github.com/deepaksinghkushwah/shop-microservices/pkg/response"
	"github.com/deepaksinghkushwah/shop-microservices/services/order-service/internal/model"
	service "github.com/deepaksinghkushwah/shop-microservices/services/order-service/internal/serivce"
	"github.com/gin-gonic/gin"
)

type CreateOrderItemRequest struct {
	OrderID      uint    `json:"order_id" binding:"required"`
	ProductID    uint    `json:"product_id" binding:"required"`
	ProductName  string  `json:"product_name" binding:"required"`
	ProductSlug  string  `json:"product_slug" binding:"required"`
	VariantID    uint    `json:"variant_id" binding:"required"`
	VariantSKU   string  `json:"variant_sku" binding:"required"`
	VariantPrice float64 `json:"variant_price" binding:"required"`
	Quantity     int     `json:"quantity" binding:"required,gt=0"`
	TotalPrice   float64 `json:"total_price" binding:"required"`
}

type UpdateOrderItemRequest struct {
	Quantity   int     `json:"quantity" binding:"required,gt=0"`
	TotalPrice float64 `json:"total_price" binding:"required"`
}

// CreateOrderItem godoc
// @Summary Create order item
// @Description Add a new item to an order
// @Tags order-items
// @Accept json
// @Produce json
// @Param item body CreateOrderItemRequest true "Order Item Data"
// @Success 201 {object} model.OrderItem
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /order-items [post]
func CreateOrderItem(c *gin.Context) {
	var req CreateOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err.Error())
		return
	}

	orderItem := &model.OrderItem{
		OrderID:      req.OrderID,
		ProductID:    req.ProductID,
		ProductName:  req.ProductName,
		ProductSlug:  req.ProductSlug,
		VariantID:    req.VariantID,
		VariantSKU:   req.VariantSKU,
		VariantPrice: req.VariantPrice,
		Quantity:     req.Quantity,
		TotalPrice:   req.TotalPrice,
	}

	createdItem, err := service.CreateOrderItem(orderItem)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, createdItem)
}

// GetOrderItem godoc
// @Summary Get order item
// @Description Get a specific order item by ID
// @Tags order-items
// @Produce json
// @Param id path uint true "Order Item ID"
// @Success 200 {object} model.OrderItem
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /order-items/:id [get]
func GetOrderItem(c *gin.Context) {
	itemID := c.GetUint("id")

	orderItem, err := service.GetOrderItem(itemID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, orderItem)
}

// ListOrderItems godoc
// @Summary List order items
// @Description Get all items for a specific order
// @Tags order-items
// @Produce json
// @Param order_id path uint true "Order ID"
// @Success 200 {object} []model.OrderItem
// @Failure 500 {object} map[string]interface{}
// @Router /orders/:order_id/items [get]
func ListOrderItems(c *gin.Context) {
	orderID := c.GetUint("order_id")

	orderItems, err := service.GetOrderItemsByOrder(orderID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, orderItems)
}

// UpdateOrderItem godoc
// @Summary Update order item
// @Description Update quantity and total price of an order item
// @Tags order-items
// @Accept json
// @Produce json
// @Param id path uint true "Order Item ID"
// @Param item body UpdateOrderItemRequest true "Updated Item Data"
// @Success 200 {object} model.OrderItem
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /order-items/:id [put]
func UpdateOrderItem(c *gin.Context) {
	var req UpdateOrderItemRequest
	itemID := c.GetUint("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err.Error())
		return
	}

	orderItem, err := service.UpdateOrderItem(itemID, req.Quantity, req.TotalPrice)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, orderItem)
}

// DeleteOrderItem godoc
// @Summary Delete order item
// @Description Delete an item from an order
// @Tags order-items
// @Produce json
// @Param id path uint true "Order Item ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /order-items/:id [delete]
func DeleteOrderItem(c *gin.Context) {
	itemID := c.GetUint("id")

	err := service.DeleteOrderItem(itemID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "order item deleted successfully"})
}
