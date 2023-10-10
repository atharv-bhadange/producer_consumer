package routes

import "github.com/gofiber/fiber/v2"

func InitRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server is up and running")
	})

	UserRoutes(api)
	ProductRoutes(api)
}
