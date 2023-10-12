package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atharv-bhadange/producer_consumer/database"
	"github.com/atharv-bhadange/producer_consumer/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testDB() (*gorm.DB, error) {
	gormDB, err := gorm.Open(sqlite.Dialector{}, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	database.DB.Db = gormDB

	return gormDB, nil
}

func TestGetUser(t *testing.T) {

	// Create gorm database
	gormDB, err := testDB()

	if err != nil {
		t.Fatalf("Error creating gorm database: %v", err)
	}
	// Create fiber app
	app := fiber.New()

	// Create database instance

	// Create mock data
	mockUser := []map[string]interface{}{
		{
			"id":         1,
			"name":       "purva",
			"mobile":     "2234567890",
			"latitude":   123.123,
			"longitude":  123.321,
			"created_at": "2023-10-10T11:20:18.728189Z",
			"updated_at": "2023-10-10T11:20:18.728189Z",
		},
	}

	// create users table
	if err := gormDB.Table("users").AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("Error creating mock table: %v", err)
	}

	// Add mock data to users table
	if err := gormDB.Table("users").Create(&mockUser).Error; err != nil {
		t.Fatalf("Error adding mock data to table: %v", err)
	}

	// Add route to fiber app
	app.Get("/user/:id", GetUser)

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Get existing user",
			id:             "1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":200,"message":"User found","data":{"user":{"id":1,"user_name":"purva","user_mobile":"2234567890","user_latitude":123.123,"user_longitude":123.321,"created_at":"2023-10-10T11:20:18.728189Z","updated_at":"2023-10-10T11:20:18.728189Z"}}}`,
		},
		{
			name:           "Get non-existing user",
			id:             "100",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":404,"message":"User not found","data":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user/%s", tt.id), nil)
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

			fmt.Println(string(body))
			fmt.Println(tt.expectedBody)

			if string(body) != tt.expectedBody {
				t.Errorf("Expected response body %s, but got %s", tt.expectedBody, string(body))
			}
		})
	}
}
