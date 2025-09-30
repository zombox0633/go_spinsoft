package config

import (
	"log"
	"time"

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
				"error":     true,
				"message":   message,
				"timestamp": time.Now().Format(time.RFC3339),
			})
		},
	})

	application := &ApplicationType{
		fiber:  app,
		config: cfg,
	}

	middleware.SetupCorsMiddleware(app)

	application.fiber.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	return application
}

func (app *ApplicationType) Start() error {
	return app.fiber.Listen(":" + app.config.Port)
}

func (app *ApplicationType) Shutdown() error {
	log.Println("Gracefully shutting down Fiber server...")
	return app.fiber.Shutdown()
}
