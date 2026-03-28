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

	authPkg "github.com/deepaksinghkushwah/shop-microservices/pkg/auth"
	"github.com/deepaksinghkushwah/shop-microservices/services/auth-service/internal/model"
	"github.com/deepaksinghkushwah/shop-microservices/services/auth-service/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupAuthTestApp(t *testing.T) *gin.Engine {
	// Setup PostgreSQL test container
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "auth_test",
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
	dsn := fmt.Sprintf("host=%s port=%s user=postgres password=password dbname=auth_test sslmode=disable",
		host, port.Port())

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Use this DB for the app
	database.DB = db

	// Run migrations
	err = database.DB.AutoMigrate(&model.User{})
	if err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	// Register routes
	r.POST("/register", Register)
	r.POST("/login", Login)

	auth := r.Group("/")
	auth.Use(authPkg.AuthMiddleware())
	auth.GET("/profile", Profile)

	return r
}

func mustAuthJSON(t *testing.T, body interface{}) []byte {
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}
	return b
}

func doAuthRequest(t *testing.T, client *http.Client, method, url string, body interface{}) *http.Response {
	var buf io.Reader
	if body != nil {
		buf = bytes.NewReader(mustAuthJSON(t, body))
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

func readAuthBody(t *testing.T, resp *http.Response) []byte {
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	return b
}

// TestAuthService_Register tests user registration
func TestAuthService_Register(t *testing.T) {
	r := setupAuthTestApp(t)
	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	// Test successful registration
	regReq := map[string]interface{}{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "password123",
	}
	resp := doAuthRequest(t, client, http.MethodPost, server.URL+"/register", regReq)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body := readAuthBody(t, resp)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("unable to parse response json: %v", err)
	}

	if _, ok := result["message"]; !ok {
		t.Fatalf("response missing message field")
	}
}

// TestAuthService_Login tests user login and JWT token generation
func TestAuthService_Login(t *testing.T) {
	r := setupAuthTestApp(t)
	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	// First register a user
	regReq := map[string]interface{}{
		"name":     "Jane Doe",
		"email":    "jane@example.com",
		"password": "password456",
	}
	resp := doAuthRequest(t, client, http.MethodPost, server.URL+"/register", regReq)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("registration failed: expected 200, got %d", resp.StatusCode)
	}

	// Now login
	loginReq := map[string]interface{}{
		"email":    "jane@example.com",
		"password": "password456",
	}
	resp = doAuthRequest(t, client, http.MethodPost, server.URL+"/login", loginReq)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login failed: expected 200, got %d", resp.StatusCode)
	}

	body := readAuthBody(t, resp)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("unable to parse response json: %v", err)
	}

	if _, ok := result["token"]; !ok {
		t.Fatalf("response missing token field")
	}
}

// TestAuthService_LoginInvalidCredentials tests login with invalid credentials
func TestAuthService_LoginInvalidCredentials(t *testing.T) {
	r := setupAuthTestApp(t)
	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	// Try login with non-existent user
	loginReq := map[string]interface{}{
		"email":    "nonexistent@example.com",
		"password": "wrongpassword",
	}
	resp := doAuthRequest(t, client, http.MethodPost, server.URL+"/login", loginReq)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid credentials, got %d", resp.StatusCode)
	}
}

// TestAuthService_RegisterAndLoginWorkflow tests complete auth flow
func TestAuthService_RegisterAndLoginWorkflow(t *testing.T) {
	r := setupAuthTestApp(t)
	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	// Step 1: Register user
	regReq := map[string]interface{}{
		"name":     "Bob Smith",
		"email":    "bob@example.com",
		"password": "securepass123",
	}
	resp := doAuthRequest(t, client, http.MethodPost, server.URL+"/register", regReq)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("registration failed: expected 200, got %d", resp.StatusCode)
	}

	// Step 2: Login user
	loginReq := map[string]interface{}{
		"email":    "bob@example.com",
		"password": "securepass123",
	}
	resp = doAuthRequest(t, client, http.MethodPost, server.URL+"/login", loginReq)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login failed: expected 200, got %d", resp.StatusCode)
	}

	body := readAuthBody(t, resp)
	var loginResult map[string]interface{}
	if err := json.Unmarshal(body, &loginResult); err != nil {
		t.Fatalf("unable to parse login response: %v", err)
	}

	token, ok := loginResult["token"]
	if !ok {
		t.Fatalf("login response missing token field")
	}

	// Step 3: Access protected profile endpoint with token
	profileReq, err := http.NewRequest(http.MethodGet, server.URL+"/profile", nil)
	if err != nil {
		t.Fatalf("failed to create profile request: %v", err)
	}
	profileReq.Header.Set("Authorization", "Bearer "+token.(string))

	resp, err = client.Do(profileReq)
	if err != nil {
		t.Fatalf("profile request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body = readAuthBody(t, resp)
	var profileResult map[string]interface{}
	if err := json.Unmarshal(body, &profileResult); err != nil {
		t.Fatalf("unable to parse profile response: %v", err)
	}

	// Verify profile response contains data
	if _, ok := profileResult["data"]; !ok {
		t.Fatalf("profile response missing data field")
	}
}
