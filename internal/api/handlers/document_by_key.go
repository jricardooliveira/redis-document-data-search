package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jricardooliveira/redis-document-data-search/internal/redisutil"
)

// DocumentByKeyHandler returns the raw JSON document for a given Redis key (customer:event:...)
func DocumentByKeyHandler(redisURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		key := c.Query("key")
		if key == "" {
			return c.Status(400).JSON(fiber.Map{
				"error":         "missing key parameter",
				"query_time_ms": time.Since(start).Milliseconds(),
			})
		}
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":         "internal server error: redis client",
				"query_time_ms": time.Since(start).Milliseconds(),
			})
		}
		ctx := context.Background()
		jsonStr, err := client.Do(ctx, "JSON.GET", key, "$").Text()
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error":         "not found",
				"query_time_ms": time.Since(start).Milliseconds(),
			})
		}
		var doc interface{}
if err := json.Unmarshal([]byte(jsonStr), &doc); err != nil {
	return c.Status(500).JSON(fiber.Map{
		"error":         "failed to parse JSON from Redis",
		"query_time_ms": time.Since(start).Milliseconds(),
	})
}
return c.JSON(fiber.Map{
	"key":           key,
	"document":      doc,
	"query_time_ms": time.Since(start).Milliseconds(),
})
	}
}
