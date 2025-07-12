package handlers

import (
	"context"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
	"runtime"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/jricardooliveira/redis-document-data-search/internal/faker"
	"github.com/jricardooliveira/redis-document-data-search/internal/redisutil"
)

func GenerateEventsHandler(redisURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		count, _ := strconv.Atoi(c.Query("count", "1000"))
		start := time.Now()
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "internal server error: redis client"})
		}

		concurrency := runtime.NumCPU() / 2
		if concurrency < 1 {
			concurrency = 1
		}
		if count < concurrency {
			concurrency = count
		}
		var wg sync.WaitGroup
		sem := make(chan struct{}, concurrency)
		errCh := make(chan error, count)

		for i := 0; i < count; i++ {
			wg.Add(1)
			sem <- struct{}{}
			go func(i int) {
				defer wg.Done()
				event := faker.RandomEvent()
				key := "event:" + strconv.Itoa(i)
				if err := redisutil.StoreJSON(client, key, event); err != nil {
					errCh <- err
				}
				<-sem
			}(i)
		}
		wg.Wait()
		close(errCh)

		if len(errCh) > 0 {
			queryTimeMs := time.Since(start).Milliseconds()
			return c.Status(500).JSON(fiber.Map{"error": (<-errCh).Error(), "query_time_ms": queryTimeMs})
		}
		queryTimeMs := time.Since(start).Milliseconds()
		return PrettyJSON(c, fiber.Map{"status": "ok", "stored": count, "query_time_ms": queryTimeMs})
	}
}

func SearchEventsHandler(redisURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "internal server error: redis client"})
		}
		identifiers := map[string]string{}
		for k, v := range c.Queries() {
			if k != "limit" && k != "offset" {
				identifiers[k] = v
			}
		}
		limit, err1 := strconv.Atoi(c.Query("limit", "10"))
		if err1 != nil || limit < 1 {
			return c.Status(400).JSON(fiber.Map{"error": "limit must be a positive integer"})
		}
		offset, err2 := strconv.Atoi(c.Query("offset", "0"))
		if err2 != nil || offset < 0 {
			return c.Status(400).JSON(fiber.Map{"error": "offset must be a non-negative integer"})
		}
		query := BuildRediSearchQuery(identifiers)
		start := time.Now()
		results, err := redisutil.SearchFTSWithLimit(client, "eventIdx", query, limit, offset)
		queryTimeMs := time.Since(start).Milliseconds()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error(), "query_time_ms": queryTimeMs})
		}
		return PrettyJSON(c, fiber.Map{"results": results, "query_time_ms": queryTimeMs})
	}
}

func RandomEventHandler(redisURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		ctx := context.Background()
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			queryTimeMs := time.Since(start).Milliseconds()
			return c.Status(500).JSON(fiber.Map{"error": "internal server error: redis client", "query_time_ms": queryTimeMs})
		}
		iter := client.Scan(ctx, 0, "event:*", 1000).Iterator()
		var keys []string
		for iter.Next(ctx) {
			keys = append(keys, iter.Val())
		}
		if len(keys) == 0 {
			queryTimeMs := time.Since(start).Milliseconds()
			return c.Status(404).JSON(fiber.Map{"error": "no events found", "query_time_ms": queryTimeMs})
		}
		key := keys[rand.Intn(len(keys))]
		val, err := client.Do(ctx, "JSON.GET", key, "$" ).Text()
		if err != nil {
			queryTimeMs := time.Since(start).Milliseconds()
			return c.Status(500).JSON(fiber.Map{"error": err.Error(), "query_time_ms": queryTimeMs})
		}
		var arr []interface{}
		queryTimeMs := time.Since(start).Milliseconds()
		if err := json.Unmarshal([]byte(val), &arr); err == nil && len(arr) > 0 {
			return PrettyJSON(c, fiber.Map{"result": arr[0], "query_time_ms": queryTimeMs})
		}
		return c.SendString(val)
	}
}

