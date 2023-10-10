package handlers

import (
	"fmt"
	"net/http"

	"github.com/atharv-bhadange/producer_consumer/database"
	"github.com/atharv-bhadange/producer_consumer/models"
	"github.com/gofiber/fiber/v2"
)

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	// print request body json
	fmt.Println(string(c.Body()))

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			models.Response{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
				Data:   nil,
			},
		)
	}
	// Create user in database use gorm

	db := database.DB.Db

	if err := db.Where("mobile = ?", user.Mobile).First(&user).Error; err == nil {
		return c.Status(http.StatusConflict).JSON(
			models.Response{
				Status:  http.StatusConflict,
				Message: "User already exists",
				Data:  nil,
			},
		)
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			models.Response{
				Status:  http.StatusInternalServerError,
				Message: "Unable to create user",
				Data:  nil,
			},
		)
	}

	return c.Status(http.StatusCreated).JSON(
		models.Response{
			Status:  http.StatusCreated,
			Message: "User created successfully",
			Data: fiber.Map{"user": user},
		},
	)
}

func GetUser(c *fiber.Ctx) error {
	var user models.User

	db := database.DB.Db

	fmt.Println(c.Params("id"))

	if err := db.Where("id = ?", c.Params("id")).First(&user).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(
			models.Response{
				Status:  http.StatusNotFound,
				Message: "User not found",
				Data:  nil,
			},
		)
	}

	return c.Status(http.StatusOK).JSON(
		models.Response{
			Status:  http.StatusOK,
			Message: "User found",
			Data: fiber.Map{"user": user},
		},
	)
}

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User

	db := database.DB.Db

	if err := db.Find(&users).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			models.Response{
				Status:  http.StatusInternalServerError,
				Message: "Unable to fetch users",
				Data:  nil,
			},
		)
	}

	return c.Status(http.StatusOK).JSON(
		models.Response{
			Status:  http.StatusOK,
			Message: "Users found",
			Data: fiber.Map{"users": users},
		},
	)
}
