package router

import (
	handler "gateway-settings-api/controllers"
	"gateway-settings-api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Welcome to Gateway Settings API."})
	})

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)

	// Contract Allowlist
	contractAllowlist := api.Group("/contract-allowlist")
	contractAllowlist.Post("/add-contract", middleware.Protected(), handler.AddContractToAllowlist)
}
