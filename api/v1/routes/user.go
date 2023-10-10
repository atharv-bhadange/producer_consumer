package routes

import (
	"github.com/atharv-bhadange/producer_consumer/api/v1/handlers"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router fiber.Router) {
	
	router.Post("/user", handlers.CreateUser)
	router.Get("/user/:id", handlers.GetUser)
	router.Get("/user", handlers.GetAllUsers)
}
