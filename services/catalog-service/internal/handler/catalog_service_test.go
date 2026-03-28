package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/catalog-service/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupCatalogTestApp(t *testing.T) *gin.Engine {
	// Setup PostgreSQL test container
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "catalog_test",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to create test container: %v", err)
	}

	t.Cleanup(func() {
		_ = container.Terminate(ctx)
	})

	// Get the container's host and port
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	// Create connection string
	dsn := fmt.Sprintf("host=%s port=%s user=postgres password=password dbname=catalog_test sslmode=disable",
		host, port.Port())

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Use this DB for the app
	database.DB = db

	// Run migrations
	err = database.DB.AutoMigrate(
		&model.Category{},
		&model.Product{},
		&model.ProductImage{},
		&model.ProductVariant{},
		&model.Attribute{},
		&model.AttributeValue{},
		&model.VariantAttributeValue{},
		&model.ProductAttribute{},
	)
	if err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	// Register routes
	r.POST("/categories", CreateCategory)
	r.GET("/categories", ListCategories)

	r.POST("/products", CreateProduct)
	r.GET("/products", ListProducts)
	r.GET("/products/:slug", GetProduct)

	r.POST("/product-images", CreateProductImage)
	r.GET("/product-images", ListProductImages)

	r.POST("/variants", CreateVariant)

	r.POST("/attributes", CreateAttribute)
	r.GET("/attributes", ListAttributes)

	r.POST("/attributes/:id/values", CreateAttributeValue)
	r.GET("/attributes/:id/values", ListAttributeValues)

	r.POST("/products/:id/attributes", AssignAttributeToProduct)

	return r
}

func mustCatalogJSON(t *testing.T, body interface{}) []byte {
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}
	return b
}

func doCatalogRequest(t *testing.T, client *http.Client, method, url string, body interface{}) *http.Response {
	var buf io.Reader
	if body != nil {
		buf = bytes.NewReader(mustCatalogJSON(t, body))
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}

func readCatalogBody(t *testing.T, resp *http.Response) []byte {
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	return b
}

// TestCatalogService_CreateCategory tests category creation
func TestCatalogService_CreateCategory(t *testing.T) {
	r := setupCatalogTestApp(t)
	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	catReq := map[string]interface{}{
		"name": "Electronics",
		"slug": "electronics",
	}
	resp := doCatalogRequest(t, client, http.MethodPost, server.URL+"/categories", catReq)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body := readCatalogBody(t, resp)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if result["success"] != true {
		t.Fatalf("expected success=true")
	}
}

// TestCatalogService_ListCategories tests category listing
func TestCatalogService_ListCategories(t *testing.T) {
	r := setupCatalogTestApp(t)
	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	// Create a category first
	catReq := map[string]interface{}{
		"name": "Fashion",
		"slug": "fashion",
	}
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/categories", catReq)

	// List categories
	resp := doCatalogRequest(t, client, http.MethodGet, server.URL+"/categories", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body := readCatalogBody(t, resp)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if result["success"] != true {
		t.Fatalf("expected success=true")
	}
}

// TestCatalogService_CreateProduct tests product creation
func TestCatalogService_CreateProduct(t *testing.T) {
	r := setupCatalogTestApp(t)
	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	// Create category first
	catReq := map[string]interface{}{
		"name": "Books",
		"slug": "books",
	}
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/categories", catReq)

	// Create product
	prodReq := map[string]interface{}{
		"name":        "Go Programming",
		"slug":        "go-programming",
		"description": "Learn Go",
		"price":       29.99,
		"category_id": 1,
	}
	resp := doCatalogRequest(t, client, http.MethodPost, server.URL+"/products", prodReq)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body := readCatalogBody(t, resp)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if result["success"] != true {
		t.Fatalf("expected success=true")
	}
}

// TestCatalogService_ListProducts tests product listing with pagination
func TestCatalogService_ListProducts(t *testing.T) {
	r := setupCatalogTestApp(t)
	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	// Create category
	catReq := map[string]interface{}{
		"name": "Software",
		"slug": "software",
	}
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/categories", catReq)

	// Create product
	prodReq := map[string]interface{}{
		"name":        "Python Course",
		"slug":        "python-course",
		"description": "Learn Python",
		"price":       49.99,
		"category_id": 1,
	}
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/products", prodReq)

	// List products
	resp := doCatalogRequest(t, client, http.MethodGet, server.URL+"/products?page=1&limit=10", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body := readCatalogBody(t, resp)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	data := result["data"].(map[string]interface{})
	products := data["products"].([]interface{})
	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}
}

// TestCatalogService_GetProductBySlug tests fetching a product by slug
func TestCatalogService_GetProductBySlug(t *testing.T) {
	r := setupCatalogTestApp(t)
	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	// Create category and product
	catReq := map[string]interface{}{"name": "Home", "slug": "home"}
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/categories", catReq)

	prodReq := map[string]interface{}{
		"name":        "Smart Bulb",
		"slug":        "smart-bulb",
		"description": "IoT bulb",
		"price":       19.99,
		"category_id": 1,
	}
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/products", prodReq)

	// Get product by slug
	resp := doCatalogRequest(t, client, http.MethodGet, server.URL+"/products/smart-bulb", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body := readCatalogBody(t, resp)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if result["success"] != true {
		t.Fatalf("expected success=true")
	}
}

// TestCatalogService_CreateProductImage tests product image creation
func TestCatalogService_CreateProductImage(t *testing.T) {
	r := setupCatalogTestApp(t)
	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	// Create category and product
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/categories", map[string]interface{}{"name": "Auto", "slug": "auto"})
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/products", map[string]interface{}{
		"name": "Car", "slug": "car", "description": "Nice car", "price": 25000.0, "category_id": 1,
	})

	// Create image
	imgReq := map[string]interface{}{
		"product_id": 1,
		"url":        "https://example.com/car.jpg",
		"is_primary": true,
	}
	resp := doCatalogRequest(t, client, http.MethodPost, server.URL+"/product-images", imgReq)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body := readCatalogBody(t, resp)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if result["success"] != true {
		t.Fatalf("expected success=true")
	}
}

// TestCatalogService_CreateAttribute tests attribute creation
func TestCatalogService_CreateAttribute(t *testing.T) {
	r := setupCatalogTestApp(t)
	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	attrReq := map[string]interface{}{
		"name": "Color",
		"slug": "color",
	}
	resp := doCatalogRequest(t, client, http.MethodPost, server.URL+"/attributes", attrReq)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	body := readCatalogBody(t, resp)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	// Attribute handler returns the attribute directly, not wrapped
	t.Logf("CreateAttribute response: %v", result)
	hasID := false
	hasName := false
	if _, ok := result["id"]; ok {
		hasID = true
	} else if _, ok := result["ID"]; ok {
		hasID = true
	}
	if _, ok := result["name"]; ok {
		hasName = true
	} else if _, ok := result["Name"]; ok {
		hasName = true
	}
	if !hasID && !hasName {
		t.Fatalf("response missing attribute fields: %v", result)
	}
}

// TestCatalogService_CreateAttributeValue tests attribute value creation
func TestCatalogService_CreateAttributeValue(t *testing.T) {
	r := setupCatalogTestApp(t)
	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	// Create attribute first
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/attributes", map[string]interface{}{"name": "Size", "slug": "size"})

	// Create attribute value
	valReq := map[string]interface{}{
		"value": "Large",
	}
	resp := doCatalogRequest(t, client, http.MethodPost, server.URL+"/attributes/1/values", valReq)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	body := readCatalogBody(t, resp)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	// AttributeValue handler returns the value directly
	t.Logf("CreateAttributeValue response: %v", result)
	hasID := false
	hasValue := false
	if _, ok := result["id"]; ok {
		hasID = true
	} else if _, ok := result["ID"]; ok {
		hasID = true
	}
	if _, ok := result["value"]; ok {
		hasValue = true
	} else if _, ok := result["Value"]; ok {
		hasValue = true
	}
	if !hasID && !hasValue {
		t.Fatalf("response missing attribute value fields: %v", result)
	}
}

// TestCatalogService_ListAttributeValues tests listing attribute values
func TestCatalogService_ListAttributeValues(t *testing.T) {
	r := setupCatalogTestApp(t)
	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	// Create attribute and values
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/attributes", map[string]interface{}{"name": "Material", "slug": "material"})
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/attributes/1/values", map[string]interface{}{"value": "Cotton"})
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/attributes/1/values", map[string]interface{}{"value": "Polyester"})

	// List values
	resp := doCatalogRequest(t, client, http.MethodGet, server.URL+"/attributes/1/values", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body := readCatalogBody(t, resp)
	var result []map[string]interface{}
	json.Unmarshal(body, &result)
	// List returns a JSON array of attribute values
	if len(result) < 1 {
		t.Fatalf("expected at least 1 attribute value")
	}
}

// TestCatalogService_CreateVariant tests creating product variant
func TestCatalogService_CreateVariant(t *testing.T) {
	r := setupCatalogTestApp(t)
	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	// Create category and product
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/categories", map[string]interface{}{"name": "Clothing", "slug": "clothing"})
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/products", map[string]interface{}{
		"name": "T-Shirt", "slug": "t-shirt", "description": "Casual", "price": 15.0, "category_id": 1,
	})

	// Create variant
	varReq := map[string]interface{}{
		"product_id": 1,
		"sku":        "TSHIRT-001-S",
		"price":      15.99,
		"stock":      50,
	}
	resp := doCatalogRequest(t, client, http.MethodPost, server.URL+"/variants", varReq)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body := readCatalogBody(t, resp)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if result["success"] != true {
		t.Fatalf("expected success=true")
	}
}

// TestCatalogService_AssignAttributeToProduct tests assigning attribute to product
func TestCatalogService_AssignAttributeToProduct(t *testing.T) {
	r := setupCatalogTestApp(t)
	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	// Setup: category, product, attribute
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/categories", map[string]interface{}{"name": "Shoes", "slug": "shoes"})
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/products", map[string]interface{}{
		"name": "Sneakers", "slug": "sneakers", "description": "Good shoes", "price": 79.99, "category_id": 1,
	})
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/attributes", map[string]interface{}{"name": "Style", "slug": "style"})

	// Assign attribute to product
	assignReq := map[string]interface{}{
		"attribute_id": 1,
	}
	resp := doCatalogRequest(t, client, http.MethodPost, server.URL+"/products/1/attributes", assignReq)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body := readCatalogBody(t, resp)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	// AssignAttribute returns a message field
	if _, ok := result["message"]; !ok {
		t.Fatalf("expected message field in response")
	}
}

// TestCatalogService_CompleteWorkflow tests a complete catalog workflow
func TestCatalogService_CompleteWorkflow(t *testing.T) {
	r := setupCatalogTestApp(t)
	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	// 1. Create category
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/categories", map[string]interface{}{"name": "Mobile", "slug": "mobile"})

	// 2. Create product
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/products", map[string]interface{}{
		"name": "iPhone 15", "slug": "iphone-15", "description": "Latest phone", "price": 999.0, "category_id": 1,
	})

	// 3. Add product image
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/product-images", map[string]interface{}{
		"product_id": 1, "url": "https://example.com/iphone.jpg", "is_primary": true,
	})

	// 4. Create attribute
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/attributes", map[string]interface{}{"name": "Color", "slug": "color"})

	// 5. Add attribute values
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/attributes/1/values", map[string]interface{}{"value": "Black"})
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/attributes/1/values", map[string]interface{}{"value": "White"})

	// 6. Create variant
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/variants", map[string]interface{}{
		"product_id": 1, "sku": "IPHONE15-001", "price": 999.99, "stock": 100,
	})

	// 7. Assign attribute to product
	doCatalogRequest(t, client, http.MethodPost, server.URL+"/products/1/attributes", map[string]interface{}{"attribute_id": 1})

	// 8. Verify by fetching product
	resp := doCatalogRequest(t, client, http.MethodGet, server.URL+"/products/iphone-15", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body := readCatalogBody(t, resp)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if result["success"] != true {
		t.Fatalf("expected success=true in final product fetch")
	}
}
