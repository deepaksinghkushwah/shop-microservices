package handler

import (
	"net/http"

	"github.com/deepaksinghkushwah/shop-microservices/pkg/response"
	"github.com/deepaksinghkushwah/shop-microservices/services/auth-service/internal/service"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := service.Register(req.Name, req.Email, req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Registered"})
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := service.Login(req.Email, req.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Credentials",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := service.Profile(userID.(uint))
	if err != nil {
		response.Error(c, "invalid credentials")
		return
	}

	response.Success(c, user)
}
