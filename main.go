package main

import (
	"log"

	"github.com/atharv-bhadange/producer_consumer/configs"
	"github.com/atharv-bhadange/producer_consumer/database"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = database.InitDatabase()

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	port := configs.GetPort()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server is up and running")
	})

	app.Listen(port)

}
