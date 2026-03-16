package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/pkg/database"
)

type CreateAttributeRequest struct {
	Name string `json:"name" example:"Color"`
	Slug string `json:"slug" example:"color"`
}

// CreateAttribute godoc
// @Summary Create attribute
// @Description Create a new product attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param attribute body CreateAttributeRequest true "Attribute Data"
// @Success 201 {object} map[string]interface{}
// @Router /attributes [post]
func CreateAttribute(c *gin.Context) {

	var req CreateAttributeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	attr := model.Attribute{
		Name: req.Name,
		Slug: req.Slug,
	}

	database.DB.Create(&attr)

	c.JSON(201, attr)
}

// ListAttributes godoc
// @Summary List attributes
// @Description Get all product attributes with their values
// @Tags attributes
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /attributes [get]
func ListAttributes(c *gin.Context) {

	var attributes []model.Attribute

	database.DB.Preload("Values").Find(&attributes)

	c.JSON(200, attributes)
}
