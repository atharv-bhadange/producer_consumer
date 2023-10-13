package main

import (
	"log"

	"github.com/atharv-bhadange/producer_consumer/api/v1/routes"
	"github.com/atharv-bhadange/producer_consumer/configs"
	"github.com/atharv-bhadange/producer_consumer/database"
	"github.com/atharv-bhadange/producer_consumer/producer"
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

	err = producer.ConnectQueue()

	if err != nil {
		log.Println(err)
		log.Fatal("Error connecting to queue")
	}

	port := configs.GetPort()

	app := fiber.New()

	routes.InitRoutes(app)

	app.Listen(port)

}
