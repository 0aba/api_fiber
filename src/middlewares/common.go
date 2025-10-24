package middlewares

import "github.com/gofiber/fiber/v2"

func SetDefautlHeaderMiddleware(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Frame-Options", "DENY")

	return c.Next()
}

func SetHeaderV1Middleware(c *fiber.Ctx) error {
	c.Set("X-API-Version", "1.0")

	return c.Next()
}
