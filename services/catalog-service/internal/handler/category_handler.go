package handler

import (
	"github.com/deepaksinghkushwah/shop-microservices/pkg/response"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/service"
	"github.com/gin-gonic/gin"
)

type CreateCategoryRequest struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	ParentID *uint  `json:"parent_id"`
}

func CreateCategory(c *gin.Context) {

	var req CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "invalid request")
		return
	}

	category, err := service.CreateCategory(req.Name, req.Slug, req.ParentID)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, category)
}

func ListCategories(c *gin.Context) {

	categories, err := service.ListCategories()

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, categories)
}
