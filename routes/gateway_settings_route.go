package routes

import (
	"gateway-settings-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func AddContractToAllowlistRoute(app *fiber.App) {
	app.Post("/add-contract", controllers.AddContractToAllowlist)
}
