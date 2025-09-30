package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func APIKeyMiddleware(keyValue string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")

		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "API Key is required",
			})
		}

		if apiKey != keyValue {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid API Key",
			})
		}

		return c.Next()
	}
}
