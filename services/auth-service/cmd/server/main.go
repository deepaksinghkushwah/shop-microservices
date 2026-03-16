package main

import (
	authPkg "github.com/deepaksinghkushwah/shop-microservices/pkg/auth"
	"github.com/deepaksinghkushwah/shop-microservices/pkg/config"
	"github.com/deepaksinghkushwah/shop-microservices/services/auth-service/internal/handler"
	"github.com/deepaksinghkushwah/shop-microservices/services/auth-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/auth-service/pkg/database"

	"github.com/gin-gonic/gin"

	_ "github.com/deepaksinghkushwah/shop-microservices/services/auth-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Auth Service API
// @version 1.0
// @description Authentication microservice
// @host localhost:8081
// @BasePath /
func main() {
	config.LoadEnv()
	err := database.Connect()
	if err != nil {
		panic(err)
	}

	database.DB.AutoMigrate(&model.User{})

	r := gin.Default()
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	auth := r.Group("/")
	auth.Use(authPkg.AuthMiddleware())
	auth.GET("/profile", handler.Profile)

	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + config.GetEnv("AUTH_SERVICE_PORT", "8081"))
}
