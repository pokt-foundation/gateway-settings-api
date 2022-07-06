package main

import (
	"gateway-settings-api/configs"
	"gateway-settings-api/router"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.Default())

	configs.ConnectDB()

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))

	app.Listen(":3000")
}
