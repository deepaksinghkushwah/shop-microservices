package handler

import (
	"strconv"

	"github.com/deepaksinghkushwah/shop-microservices/pkg/response"
	"github.com/deepaksinghkushwah/shop-microservices/pkg/validation"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/service"
	"github.com/gin-gonic/gin"
)

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3"`
	Slug        string  `json:"slug" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	CategoryID  uint    `json:"category_id" validate:"required"`
}

// CreateProduct godoc
// @Summary Create product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body CreateProductRequest true "Product Data"
// @Success 200 {object} map[string]interface{}
// @Router /products [post]
func CreateProduct(c *gin.Context) {

	var req CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "invalid request")
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		response.Error(c, "invalid request")
		return
	}

	product, err := service.CreateProduct(
		req.Name,
		req.Slug,
		req.Description,
		req.Price,
		req.CategoryID,
	)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, product)
}

// ListProducts godoc
// @Summary List products
// @Description Get paginated product list
// @Tags products
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Router /products [get]
func ListProducts(c *gin.Context) {

	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	categoryID := uint(0)
	if cid := c.Query("category_id"); cid != "" {
		if id, err := strconv.ParseUint(cid, 10, 32); err == nil {
			categoryID = uint(id)
		}
	}

	products, total, err := service.ListProducts(page, limit, categoryID)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"products": products,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func GetProduct(c *gin.Context) {

	slug := c.Param("slug")

	product, err := service.GetProduct(slug)

	if err != nil {
		response.Error(c, "product not found")
		return
	}

	response.Success(c, product)
}
