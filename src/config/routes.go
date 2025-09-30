package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zombox0633/go_spinsoft/src/station"
)

func setRoutes(app *fiber.App) {
	DB := DB.DBName
	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	// Setup routes
	station.StationRoutes(api, DB)
}
