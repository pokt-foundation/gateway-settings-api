package main

import (
	"gateway-settings-api/configs"
	"gateway-settings-api/router"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	configs.ConnectDB()

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
