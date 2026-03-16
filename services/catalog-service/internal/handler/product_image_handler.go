package handler

import (
	"github.com/deepaksinghkushwah/shop-microservices/pkg/response"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/repository"
	"github.com/gin-gonic/gin"
)

type CreateImageRequest struct {
	ProductID uint   `json:"product_id" example:"1"`
	URL       string `json:"url" example:"https://example.com/image.jpg"`
	IsPrimary bool   `json:"is_primary" example:"true"`
}

// CreateProductImage godoc
// @Summary Create product image
// @Description Create a new product image
// @Tags product-images
// @Accept json
// @Produce json
// @Param image body CreateImageRequest true "Product Image Data"
// @Success 200 {object} map[string]interface{}
// @Router /product-images [post]
func CreateProductImage(c *gin.Context) {

	var req CreateImageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "invalid request")
		return
	}

	image := model.ProductImage{
		ProductID: req.ProductID,
		URL:       req.URL,
		IsPrimary: req.IsPrimary,
	}

	err := repository.CreateProductImage(&image)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, image)
}

// ListProductImages godoc
// @Summary List product images
// @Description Get all product images
// @Tags product-images
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /product-images [get]
func ListProductImages(c *gin.Context) {

	images, err := repository.GetAllProductImages()

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, images)
}
