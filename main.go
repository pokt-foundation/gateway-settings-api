package main

import (
	"gateway-settings-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Routes
	routes.AddContractToAllowlistRoute(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Hello World"})
	})

	app.Listen(":3000")
}
