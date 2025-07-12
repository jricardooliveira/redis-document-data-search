package handlers

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/jricardooliveira/redis-document-data-search/internal/redisutil"
)

func CreateIndexesHandler(redisURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		client := redisutil.GetSingletonRedisClient(redisURL)
		err1 := redisutil.CreateCustomerIndex(client)
		err2 := redisutil.CreateEventIndex(client)
		queryTimeMs := time.Since(start).Milliseconds()
		if err1 != nil || err2 != nil {
			return c.JSON(fiber.Map{"customerIdx": err1, "eventIdx": err2, "query_time_ms": queryTimeMs})
		}
		return c.JSON(fiber.Map{"status": "ok", "query_time_ms": queryTimeMs})
	}
}
