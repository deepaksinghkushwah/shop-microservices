package handler

import (
	"github.com/deepaksinghkushwah/shop-microservices/pkg/response"
	"github.com/deepaksinghkushwah/shop-microservices/services/order-service/internal/model"
	service "github.com/deepaksinghkushwah/shop-microservices/services/order-service/internal/serivce"
	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	UserID *uint             `json:"user_id,omitempty"`
	Note   string            `json:"note"`
	Items  []model.OrderItem `json:"items"`
}

// CreateOrder godoc
// @Summary Create order
// @Description Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Param order body CreateOrderRequest true "Order Data"
// @Success 201 {object} map[string]interface{}
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err.Error())
		return
	}

	order, err := service.CreateOrder(req.UserID, req.Note, req.Items)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, order)
}

// ListUserOrders godoc
// @Summary List user orders
// @Description List orders for a specific user
// @Tags orders
// @Accept json
// @Produce json
// @Param user_id path uint true "User ID"
// @Success 200 {object} []model.Order
// @Router /orders [get]
func ListUserOrders(c *gin.Context) {
	userID := c.GetUint("user_id") // extracted from JWT

	orders, err := service.GetOrdersByUser(userID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, orders)
}

type UpdateStatusRequest struct {
	Status string `json:"status" example:"in_progress"`
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update the status of an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Param order_id path uint true "Order ID"
// @Param status body UpdateStatusRequest true "New Status"
// @Success 201 {object} map[string]interface{}
// @Router /orders/:id/status [put]
func UpdateOrderStatus(c *gin.Context) {
	var req UpdateStatusRequest
	orderID := c.GetUint("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err.Error())
		return
	}

	err := service.UpdateOrderStatus(orderID, req.Status)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "status updated"})
}
