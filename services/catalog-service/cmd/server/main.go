package main

import (
	"github.com/deepaksinghkushwah/shop-microservices/pkg/config"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/handler"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/pkg/database"
	"github.com/gin-gonic/gin"

	_ "github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Catalog Service API
// @version 1.0
// @description Product catalog microservice
// @host localhost:8082
// @BasePath /
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
		&model.ProductVariant{},
		&model.Attribute{},
		&model.AttributeValue{},
		&model.ProductAttribute{},
		&model.VariantAttributeValue{},
	)

	r := gin.Default()

	r.POST("/categories", handler.CreateCategory)
	r.GET("/categories", handler.ListCategories)

	r.POST("/products", handler.CreateProduct)
	r.GET("/products", handler.ListProducts)
	r.GET("/products/:slug", handler.GetProduct)

	r.POST("/product-images", handler.CreateProductImage)
	r.GET("/product-images", handler.ListProductImages)

	r.POST("/variants", handler.CreateVariant)

	// ATTRIBUTE ROUTES
	r.POST("/attributes", handler.CreateAttribute)
	r.GET("/attributes", handler.ListAttributes)
	r.POST("/attributes/:id/values", handler.CreateAttributeValue)
	r.GET("/attributes/:id/values", handler.ListAttributeValues)

	// Assign attribute to product
	r.POST("/products/:id/attributes", handler.AssignAttributeToProduct)

	r.SetTrustedProxies([]string{"127.0.1"})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := config.GetEnv("CATALOG_SERVICE_PORT", "8082")

	r.Run(":" + port)
}
