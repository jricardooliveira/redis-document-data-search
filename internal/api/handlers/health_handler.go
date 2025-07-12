package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jricardooliveira/redis-document-data-search/internal/redisutil"
)


func HealthHandler(redisURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "error": "redis client failed", "redis_url": redisURL})
		}
		customerCount := 0
		eventCount := 0
		ctx := context.Background()
		custIter := client.Scan(ctx, 0, "customer:*", 1000).Iterator()
		for custIter.Next(ctx) {
			customerCount++
		}
		eventIter := client.Scan(ctx, 0, "event:*", 1000).Iterator()
		for eventIter.Next(ctx) {
			eventCount++
		}
		dbIdx := "0"
		if parts := strings.Split(redisURL, "/"); len(parts) > 1 {
			dbIdx = parts[len(parts)-1]
		}
		// Get Redis memory info
		usedBytes, usedHuman, memErr := redisutil.GetRedisMemoryInfo(client)
		memInfo := fiber.Map{}
		if memErr == nil {
			memInfo["used_memory_bytes"] = usedBytes
			memInfo["used_memory_human"] = usedHuman
		}

		indexes, idxErr := redisutil.GetIndexesAndFields()
		queryTimeMs := time.Since(start).Milliseconds()
		resp := fiber.Map{
			"status":         "ok",
			"redis_url":      redisURL,
			"db_index":       dbIdx,
			"customer_count": customerCount,
			"event_count":    eventCount,
			"redis_memory":   memInfo,
			"query_time_ms":  queryTimeMs,
		}
		if idxErr == nil {
			resp["indexes"] = indexes
		}
		return PrettyJSON(c, resp)
	}
}
