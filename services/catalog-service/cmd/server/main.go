package main

import (
	"github.com/deepaksinghkushwah/shop-microservices/pkg/config"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/handler"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadEnv()

	err := database.Connect()
	if err != nil {
		panic(err)
	}

	database.DB.AutoMigrate(
		&model.Category{},
		&model.Product{},
		&model.ProductImage{},
	)

	r := gin.Default()

	r.POST("/categories", handler.CreateCategory)
	r.GET("/categories", handler.ListCategories)

	r.POST("/products", handler.CreateProduct)
	r.GET("/products", handler.ListProducts)
	r.GET("/products/:slug", handler.GetProduct)

	r.POST("/product-images", handler.CreateProductImage)
	r.GET("/product-images", handler.ListProductImages)

	port := config.GetEnv("CATALOG_SERVICE_PORT", "8082")

	r.Run(":" + port)
}
