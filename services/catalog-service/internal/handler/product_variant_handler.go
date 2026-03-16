package handler

import (
	"github.com/deepaksinghkushwah/shop-microservices/pkg/response"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/service"
	"github.com/gin-gonic/gin"
)

type CreateVariantRequest struct {
	ProductID uint    `json:"product_id" example:"1"`
	SKU       string  `json:"sku" example:"PROD-001-S"`
	Price     float64 `json:"price" example:"99.99"`
	Stock     int     `json:"stock" example:"100"`

	AttributeValueIDs []uint `json:"attribute_value_ids" example:"1,4"`
}

// CreateVariant godoc
// @Summary Create product variant
// @Description Create a new product variant
// @Tags variants
// @Accept json
// @Produce json
// @Param variant body CreateVariantRequest true "Product Variant Data"
// @Success 200 {object} map[string]interface{}
// @Router /variants [post]
func CreateVariant(c *gin.Context) {

	var req CreateVariantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "invalid request")
		return
	}

	variant, err := service.CreateVariant(
		req.ProductID,
		req.SKU,
		req.Price,
		req.Stock,
		req.AttributeValueIDs,
	)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, variant)
}
