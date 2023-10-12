package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/atharv-bhadange/producer_consumer/models"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

func TestGetProduct(t *testing.T) {
	// Create gorm database
	gormDB, err := testDB()

	if err != nil {
		t.Fatalf("Error creating gorm database: %v", err)
	}

	// Create fiber app
	app := fiber.New()

	// Add mock data to products table
	mockProduct := map[string]interface{}{
		"product_id":                1,
		"product_name":              "Test Product",
		"product_description":       "Test Product Description",
		"product_images":            pq.StringArray{"https://example.com/image1.jpg", "https://example.com/image2.jpg"},
		"product_price":             10.99,
		"compressed_product_images": pq.StringArray{"https://example.com/image1.jpg", "https://example.com/image2.jpg"},
		"created_at":                "2023-10-10T11:20:18.728189Z",
		"updated_at":                "2023-10-10T11:20:18.728189Z",
	}

	// create products table
	if err := gormDB.Table("products").AutoMigrate(&models.Product{}); err != nil {
		t.Fatalf("Error creating mock table: %v", err)
	}

	if err := gormDB.Table("products").Create(&mockProduct).Error; err != nil {
		t.Fatalf("Error adding mock data to table: %v", err)
	}

	// Add route to fiber app
	app.Get("/product/:id", GetProduct)

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Get existing product",
			id:             "1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":200,"message":"Product found","data":{"product":{"product_id":1,"product_name":"Test Product","product_description":"Test Product Description","product_images":["https://example.com/image1.jpg","https://example.com/image2.jpg"],"product_price":10.99,"compressed_product_images":["https://example.com/image1.jpg","https://example.com/image2.jpg"],"created_at":"2023-10-10T11:20:18.728189Z","updated_at":"2023-10-10T11:20:18.728189Z"}}}`,
		},
		{
			name:           "Get non-existing product",
			id:             "100",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":404,"message":"Product not found","data":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/product/"+tt.id, nil)
			resp, err := app.Test(req)

			if err != nil {
				t.Fatalf("Error making request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatus, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)

			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			if string(body) != tt.expectedBody {
				t.Errorf("Expected response body %s, but got %s", tt.expectedBody, string(body))
			}
		})
	}
}
func TestCreateProduct(t *testing.T) {
	// Create gorm database
	gormDB, err := testDB()

	if err != nil {
		t.Fatalf("Error creating gorm database: %v", err)
	}

	// Create fiber app
	app := fiber.New()

	// create products table
	if err := gormDB.Table("products").AutoMigrate(&models.Product{}); err != nil {
		t.Fatalf("Error creating mock table: %v", err)
	}

	// Add route to fiber app
	app.Post("/product", CreateProduct)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Create product with valid request body",
			requestBody:    `{"product_name":"Test Product","product_description":"Test Product Description","product_images":["https://example.com/image1.jpg","https://example.com/image2.jpg"],"product_price":10.99,"compressed_product_images":["https://example.com/image1.jpg","https://example.com/image2.jpg"]}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"status":201,"message":"Product created successfully","data":{"product_id":1}}`,
		},
		{
			name:           "Create product with invalid request body",
			requestBody:    `{"product_name":"Test Product","product_description":"Test Product Description","product_images":["https://example.com/image1.jpg","https://example.com/image2.jpg"],"product_price":"invalid_price","compressed_product_images":["https://example.com/image1.jpg","https://example.com/image2.jpg"]}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":400,"message":"Invalid request body","data":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			if err != nil {
				t.Fatalf("Error making request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatus, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)

			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			if string(body) != tt.expectedBody {
				t.Errorf("Expected response body %s, but got %s", tt.expectedBody, string(body))
			}
		})
	}
}
func TestGetAllProducts(t *testing.T) {
	// Create gorm database
	gormDB, err := testDB()

	if err != nil {
		t.Fatalf("Error creating gorm database: %v", err)
	}

	// Create fiber app
	app := fiber.New()

	// Add mock data to products table
	mockProducts := []map[string]interface{}{
		{
			"product_id":                1,
			"product_name":              "Test Product 1",
			"product_description":       "Test Product Description 1",
			"product_images":            pq.StringArray{"https://example.com/image1.jpg", "https://example.com/image2.jpg"},
			"product_price":             10.99,
			"compressed_product_images": pq.StringArray{"https://example.com/image1.jpg", "https://example.com/image2.jpg"},
			"created_at":                "2023-10-10T11:20:18.728189Z",
			"updated_at":                "2023-10-10T11:20:18.728189Z",
		},
		{
			"product_id":                2,
			"product_name":              "Test Product 2",
			"product_description":       "Test Product Description 2",
			"product_images":            pq.StringArray{"https://example.com/image3.jpg", "https://example.com/image4.jpg"},
			"product_price":             20.99,
			"compressed_product_images": pq.StringArray{"https://example.com/image3.jpg", "https://example.com/image4.jpg"},
			"created_at":                "2023-10-10T11:20:18.728189Z",
			"updated_at":                "2023-10-10T11:20:18.728189Z",
		},
	}


	// create products table
	if err := gormDB.Table("products").AutoMigrate(&models.Product{}); err != nil {
		t.Fatalf("Error creating mock table: %v", err)
	}

	if err := gormDB.Table("products").Create(&mockProducts).Error; err != nil {
		t.Fatalf("Error adding mock data to table: %v", err)
	}

	// Add route to fiber app
	app.Get("/products", GetAllProducts)

	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Get all products",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":200,"message":"Products found","data":{"products":[{"product_id":1,"product_name":"Test Product 1","product_description":"Test Product Description 1","product_images":["https://example.com/image1.jpg","https://example.com/image2.jpg"],"product_price":10.99,"compressed_product_images":["https://example.com/image1.jpg","https://example.com/image2.jpg"],"created_at":"2023-10-10T11:20:18.728189Z","updated_at":"2023-10-10T11:20:18.728189Z"},{"product_id":2,"product_name":"Test Product 2","product_description":"Test Product Description 2","product_images":["https://example.com/image3.jpg","https://example.com/image4.jpg"],"product_price":20.99,"compressed_product_images":["https://example.com/image3.jpg","https://example.com/image4.jpg"],"created_at":"2023-10-10T11:20:18.728189Z","updated_at":"2023-10-10T11:20:18.728189Z"}]}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/products", nil)
			resp, err := app.Test(req)

			if err != nil {
				t.Fatalf("Error making request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatus, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)

			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			if string(body) != tt.expectedBody {
				t.Errorf("Expected response body %s, but got %s", tt.expectedBody, string(body))
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	// Create gorm database
	gormDB, err := testDB()

	if err != nil {
		t.Fatalf("Error creating gorm database: %v", err)
	}

	// Create fiber app
	app := fiber.New()

	// Add mock data to products table
	mockProduct := map[string]interface{}{
		"product_id":                1,
		"product_name":              "Test Product",
		"product_description":       "Test Product Description",
		"product_images":            pq.StringArray{"https://example.com/image1.jpg", "https://example.com/image2.jpg"},
		"product_price":             10.99,
		"compressed_product_images": pq.StringArray{"https://example.com/image1.jpg", "https://example.com/image2.jpg"},
		"created_at":                "2023-10-10T11:20:18.728189Z",
		"updated_at":                "2023-10-10T11:20:18.728189Z",
	}

	// create products table
	if err := gormDB.Table("products").AutoMigrate(&models.Product{}); err != nil {
		t.Fatalf("Error creating mock table: %v", err)
	}

	if err := gormDB.Table("products").Create(&mockProduct).Error; err != nil {
		t.Fatalf("Error adding mock data to table: %v", err)
	}

	// Add route to fiber app
	app.Put("/product/:id", UpdateProduct)

	tests := []struct {
		name           string
		id             string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Update existing product with valid request body",
			id:             "1",
			requestBody:    `{"product_name":"Updated Test Product","product_description":"Updated Test Product Description","product_images":["https://example.com/image1.jpg","https://example.com/image2.jpg"],"product_price":20.99,"compressed_product_images":["https://example.com/image1.jpg","https://example.com/image2.jpg"]}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":200,"message":"Product updated successfully","data":{"product_id":1,"updated":true}}`,
		},
		{
			name:           "Update non-existing product",
			id:             "100",
			requestBody:    `{"product_name":"Updated Test Product","product_description":"Updated Test Product Description","product_images":["https://example.com/image1.jpg","https://example.com/image2.jpg"],"product_price":20.99,"compressed_product_images":["https://example.com/image1.jpg","https://example.com/image2.jpg"]}`,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":404,"message":"Product not found","data":null}`,
		},
		{
			name:           "Update existing product with invalid request body",
			id:             "1",
			requestBody:    `{"product_name":"Updated Test Product","product_description":"Updated Test Product Description","product_images":["https://example.com/image1.jpg","https://example.com/image2.jpg"],"product_price":"invalid_price","compressed_product_images":["https://example.com/image1.jpg","https://example.com/image2.jpg"]}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":400,"message":"Invalid request body","data":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/product/"+tt.id, strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			if err != nil {
				t.Fatalf("Error making request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatus, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)

			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			if string(body) != tt.expectedBody {
				t.Errorf("Expected response body %s, but got %s", tt.expectedBody, string(body))
			}
		})
	}
}
