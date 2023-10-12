package handlers

import (
	"fmt"
	"net/http"

	"github.com/atharv-bhadange/producer_consumer/database"
	"github.com/atharv-bhadange/producer_consumer/models"
	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	var user models.User

	db := database.DB.Db

	fmt.Println(c.Params("id"))

	if err := db.Where("id = ?", c.Params("id")).First(&user).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(
			models.Response{
				Status:  http.StatusNotFound,
				Message: "User not found",
				Data:    nil,
			},
		)
	}

	return c.Status(http.StatusOK).JSON(
		models.Response{
			Status:  http.StatusOK,
			Message: "User found",
			Data:    fiber.Map{"user": user},
		},
	)
}

