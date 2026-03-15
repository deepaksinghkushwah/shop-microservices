package handler

import (
	"github.com/deepaksinghkushwah/shop-microservices/pkg/response"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/repository"
	"github.com/gin-gonic/gin"
)

type CreateImageRequest struct {
	ProductID uint   `json:"product_id"`
	URL       string `json:"url"`
	IsPrimary bool   `json:"is_primary"`
}

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

func ListProductImages(c *gin.Context) {

	images, err := repository.GetAllProductImages()

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, images)
}
