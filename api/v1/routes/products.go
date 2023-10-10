package routes

import (
	"github.com/atharv-bhadange/producer_consumer/api/v1/handlers"
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(router fiber.Router) {
	router.Post("/product", handlers.CreateProduct)
	router.Get("/product/:id", handlers.GetProduct)
	router.Get("/product", handlers.GetAllProducts)
	router.Put("/product/:id", handlers.UpdateProduct)
}
