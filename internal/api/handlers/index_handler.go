package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jricardooliveira/redis-document-data-search/internal/redisutil"
)

func CreateIndexesHandler(redisURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "internal server error: redis client"})
		}
		err1 := redisutil.CreateCustomerIndex(client)
		err2 := redisutil.CreateEventIndex(client)
		if err1 != nil || err2 != nil {
			return c.JSON(fiber.Map{"customerIdx": err1, "eventIdx": err2})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	}
}
