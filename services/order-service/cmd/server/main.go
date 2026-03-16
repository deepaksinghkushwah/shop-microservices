// @title Order Service API
// @version 1.0
// @description Order Service API documentation
// @host localhost:8083
// @BasePath /
// @schemes http

package main

import (
	"github.com/deepaksinghkushwah/shop-microservices/pkg/config"
	"github.com/deepaksinghkushwah/shop-microservices/services/order-service/internal/handler"
	"github.com/deepaksinghkushwah/shop-microservices/services/order-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/order-service/pkg/database"
	"github.com/gin-gonic/gin"

	_ "github.com/deepaksinghkushwah/shop-microservices/services/order-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.LoadEnv()

	err := database.Connect()
	if err != nil {
		panic(err)
	}

	database.DB.AutoMigrate(
		&model.Order{},
		&model.OrderItem{},
	)

	r := gin.Default()

	// Order routes
	r.POST("/orders", handler.CreateOrder)
	r.GET("/orders", handler.ListUserOrders)               // JWT required
	r.PUT("/orders/:id/status", handler.UpdateOrderStatus) // admin

	// OrderItem routes
	r.POST("/order-items", handler.CreateOrderItem)
	r.GET("/order-items/:id", handler.GetOrderItem)
	r.GET("/orders/:order_id/items", handler.ListOrderItems)
	r.PUT("/order-items/:id", handler.UpdateOrderItem)
	r.DELETE("/order-items/:id", handler.DeleteOrderItem)

	r.SetTrustedProxies([]string{"127.0.1"})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := config.GetEnv("ORDER_SERVICE_PORT", "8083")

	r.Run(":" + port)
}
