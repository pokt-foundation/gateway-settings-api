package router

import (
	handler "gateway-settings-api/controllers"
	"gateway-settings-api/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Entrypoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Gateway Settings API.")
	})

	// Middleware
	api := app.Group("/v1")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Welcome to Gateway Settings API."})
	})

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)

	// Contract Allowlist
	contractAllowlist := api.Group("/settings")
	contractAllowlist.Post("/add-contract", middleware.Protected(), handler.AddContractToAllowlist)
}
