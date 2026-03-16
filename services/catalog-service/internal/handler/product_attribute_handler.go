package handler

import (
	"net/http"
	"strconv"

	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/repository"
	"github.com/gin-gonic/gin"
)

type AssignAttributeRequest struct {
	AttributeID uint `json:"attribute_id" example:"1"`
}

// AssignAttributeToProduct godoc
// @Summary Assign attribute to product
// @Description Assign an existing attribute to a product
// @Tags product-attributes
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param attribute body AssignAttributeRequest true "Attribute Assignment Data"
// @Success 200 {object} map[string]string
// @Router /products/{id}/attributes [post]
func AssignAttributeToProduct(c *gin.Context) {

	productIDParam := c.Param("id")

	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req AssignAttributeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = repository.AssignAttributeToProduct(uint(productID), req.AttributeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "attribute assigned to product",
	})
}
