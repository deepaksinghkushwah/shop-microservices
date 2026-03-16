package handler

import (
	"github.com/deepaksinghkushwah/shop-microservices/pkg/response"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/service"
	"github.com/gin-gonic/gin"
)

type CreateCategoryRequest struct {
	Name     string `json:"name" example:"Electronics"`
	Slug     string `json:"slug" example:"electronics"`
	ParentID *uint  `json:"parent_id" example:"1"`
}

// CreateCategory godoc
// @Summary Create category
// @Description Create a new product category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body CreateCategoryRequest true "Category Data"
// @Success 200 {object} map[string]interface{}
// @Router /categories [post]
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

// ListCategories godoc
// @Summary List categories
// @Description Get all product categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /categories [get]
func ListCategories(c *gin.Context) {

	categories, err := service.ListCategories()

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, categories)
}
