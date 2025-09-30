package config

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/zombox0633/go_spinsoft/src/middleware"
)

type ApplicationType struct {
	fiber  *fiber.App
	config *ConfigType
}

func NewApplication(cfg *ConfigType) *ApplicationType {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}

			log.Printf("Error: %v", err)

			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": message,
			})
		},
	})

	application := &ApplicationType{
		fiber:  app,
		config: cfg,
	}

	middleware.SetupCorsMiddleware(app)

	setRoutes(app)

	return application
}

func (app *ApplicationType) Start() error {
	return app.fiber.Listen(":" + app.config.Port)
}

func (app *ApplicationType) Shutdown() error {
	log.Println("Gracefully shutting down Fiber server...")
	return app.fiber.Shutdown()
}
