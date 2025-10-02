package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func APIKeyMiddleware(keyValue string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")

		if apiKey == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "API Key is required")
		}

		if apiKey != keyValue {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid API Key")
		}

		return c.Next()
	}
}
