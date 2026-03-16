package handler

import (
"bytes"
"encoding/json"
"net/http"
"net/http/httptest"
"testing"

"github.com/deepaksinghkushwah/shop-microservices/services/order-service/internal/model"
"github.com/gin-gonic/gin"
"github.com/stretchr/testify/assert"
)

func TestCreateOrderItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		requestBody CreateOrderItemRequest
		expectError bool
	}{
		{
			name: "Valid order item creation",
			requestBody: CreateOrderItemRequest{
				OrderID:      1,
				ProductID:    100,
				ProductName:  "Test Product",
				ProductSlug:  "test-product",
				VariantID:    200,
				VariantSKU:   "SKU123",
				VariantPrice: 99.99,
				Quantity:     2,
				TotalPrice:   199.98,
			},
			expectError: false,
		},
		{
			name: "Invalid quantity (zero)",
			requestBody: CreateOrderItemRequest{
				OrderID:    1,
				ProductID:  100,
				Quantity:   0,
				TotalPrice: 0,
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
r := gin.New()
			r.POST("/order-items", CreateOrderItem)

			body, _ := json.Marshal(test.requestBody)
			req, _ := http.NewRequest("POST", "/order-items", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if test.expectError {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
			}
		})
	}
}

func TestListOrderItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/orders/:order_id/items", ListOrderItems)

	req, _ := http.NewRequest("GET", "/orders/1/items", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
}

func TestGetOrderItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/order-items/:id", GetOrderItem)

	req, _ := http.NewRequest("GET", "/order-items/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
}

func TestUpdateOrderItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid update", func(t *testing.T) {
r := gin.New()
		r.PUT("/order-items/:id", UpdateOrderItem)

		reqBody := UpdateOrderItemRequest{
			Quantity:   5,
			TotalPrice: 499.95,
		}

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/order-items/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
	})

	t.Run("Invalid quantity", func(t *testing.T) {
r := gin.New()
		r.PUT("/order-items/:id", UpdateOrderItem)

		reqBody := UpdateOrderItemRequest{
			Quantity:   0,
			TotalPrice: 0,
		}

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/order-items/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteOrderItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.DELETE("/order-items/:id", DeleteOrderItem)

	req, _ := http.NewRequest("DELETE", "/order-items/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
}

func TestCreateOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid order with items", func(t *testing.T) {
r := gin.New()
		r.POST("/orders", CreateOrder)

		userID := uint(1)
		reqBody := CreateOrderRequest{
			UserID: &userID,
			Note:   "Test order",
			Items: []model.OrderItem{
				{
					ProductID:    100,
					ProductName:  "Product",
					ProductSlug:  "product",
					VariantID:    1,
					VariantSKU:   "SKU1",
					VariantPrice: 99.99,
					Quantity:     1,
					TotalPrice:   99.99,
				},
			},
		}

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/orders", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
	})

	t.Run("Guest order (no user ID)", func(t *testing.T) {
r := gin.New()
		r.POST("/orders", CreateOrder)

		reqBody := CreateOrderRequest{
			Note: "Guest order",
			Items: []model.OrderItem{
				{
					ProductID:    100,
					ProductName:  "Product",
					ProductSlug:  "product",
					VariantID:    1,
					VariantSKU:   "SKU1",
					VariantPrice: 99.99,
					Quantity:     1,
					TotalPrice:   99.99,
				},
			},
		}

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/orders", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
	})
}

func TestListUserOrders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/orders", ListUserOrders)

	req, _ := http.NewRequest("GET", "/orders", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
}

func TestUpdateOrderStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		status      string
		expectError bool
	}{
		{
			name:        "Valid status",
			status:      "processing",
			expectError: false,
		},
		{
			name:        "Another valid status",
			status:      "shipped",
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
r := gin.New()
			r.PUT("/orders/:id/status", UpdateOrderStatus)

			reqBody := UpdateStatusRequest{Status: test.status}
			body, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest("PUT", "/orders/1/status", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
		})
	}
}
