package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/pkg/database"
)

type CreateAttributeValueRequest struct {
	Value string `json:"value" example:"Red"`
}

// CreateAttributeValue godoc
// @Summary Create attribute value
// @Description Create a new value for a specific attribute
// @Tags attribute-values
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param value body CreateAttributeValueRequest true "Attribute Value Data"
// @Success 201 {object} map[string]interface{}
// @Router /attributes/{id}/values [post]
func CreateAttributeValue(c *gin.Context) {

	attrID := c.Param("id")

	id, err := strconv.Atoi(attrID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid attribute id"})
		return
	}

	var req CreateAttributeValueRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	value := model.AttributeValue{
		AttributeID: uint(id),
		Value:       req.Value,
	}

	database.DB.Create(&value)

	c.JSON(201, value)
}

// ListAttributeValues godoc
// @Summary List attribute values
// @Description Get all values for a specific attribute
// @Tags attribute-values
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 200 {object} map[string]interface{}
// @Router /attributes/{id}/values [get]
func ListAttributeValues(c *gin.Context) {

	attrID := c.Param("id")

	id, err := strconv.Atoi(attrID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid attribute id"})
		return
	}

	var values []model.AttributeValue

	database.DB.
		Where("attribute_id = ?", id).
		Find(&values)

	c.JSON(200, values)
}
