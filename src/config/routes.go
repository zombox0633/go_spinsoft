package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zombox0633/go_spinsoft/src/middleware"
	"github.com/zombox0633/go_spinsoft/src/station"
)

func setRoutes(app *fiber.App, apiKey string) {
	if DB == nil || DB.DBName == nil {
		panic("Database not initialized")
	}

	database := DB.DBName
	api := app.Group("/api")

	api.Use(middleware.APIKeyMiddleware(apiKey))

	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello World! 😺",
		})
	})

	// Setup routes
	station.StationRoutes(api, database)
}
