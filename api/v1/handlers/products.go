package handlers

import (
	"github.com/atharv-bhadange/producer_consumer/database"
	"github.com/atharv-bhadange/producer_consumer/models"
	"github.com/atharv-bhadange/producer_consumer/producer"
	"github.com/gofiber/fiber/v2"
)

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.Response{
				Status:  fiber.StatusBadRequest,
				Message: "Invalid request body",
				Data:    nil,
			},
		)
	}

	db := database.DB.Db

	if err := db.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			models.Response{
				Status:  fiber.StatusInternalServerError,
				Message: "Unable to create product",
				Data:    nil,
			},
		)
	}

	err := producer.PublishMessage(product.ProductID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			models.Response{
				Status:  fiber.StatusInternalServerError,
				Message: "Unable to publish message",
				Data:    nil,
			},
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		models.Response{
			Status:  fiber.StatusCreated,
			Message: "Product created successfully",
			Data: fiber.Map{
				"product_id": product.ProductID,
			},
		},
	)
}

func GetProduct(c *fiber.Ctx) error {
	var product models.Product

	// Get product from database use gorm

	db := database.DB.Db

	if err := db.Where("product_id = ?", c.Params("id")).First(&product).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			models.Response{
				Status:  fiber.StatusNotFound,
				Message: "Product not found",
				Data:    nil,
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		models.Response{
			Status:  fiber.StatusOK,
			Message: "Product found",
			Data:    fiber.Map{"product": product},
		},
	)
}

func GetAllProducts(c *fiber.Ctx) error {

	var products []models.Product

	// Get all products from database use gorm
	db := database.DB.Db

	if err := db.Find(&products).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			models.Response{
				Status:  fiber.StatusNotFound,
				Message: "Products not found",
				Data:    nil,
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		models.Response{
			Status:  fiber.StatusOK,
			Message: "Products found",
			Data:    fiber.Map{"products": products},
		},
	)
}

func UpdateProduct(c *fiber.Ctx) error {

	var product models.Product

	// Update product in database use gorm
	db := database.DB.Db

	if err := db.Where("product_id = ?", c.Params("id")).First(&product).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			models.Response{
				Status:  fiber.StatusNotFound,
				Message: "Product not found",
				Data:    nil,
			},
		)
	}

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.Response{
				Status:  fiber.StatusBadRequest,
				Message: "Invalid request body",
				Data:    nil,
			},
		)
	}

	if err := db.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			models.Response{
				Status:  fiber.StatusInternalServerError,
				Message: "Unable to update product",
				Data:    nil,
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		models.Response{
			Status:  fiber.StatusOK,
			Message: "Product updated successfully",
			Data: fiber.Map{
				"product_id": product.ProductID,
				"updated":    true,
			},
		},
	)

}
