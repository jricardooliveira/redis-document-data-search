package handlers

import (
	"context"
	"encoding/json"
	"math/rand"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jricardooliveira/redis-document-data-search/internal/faker"
	"github.com/jricardooliveira/redis-document-data-search/internal/redisutil"
)

func GenerateCustomersHandler(redisURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		count, _ := strconv.Atoi(c.Query("count", "1000"))
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "internal server error: redis client"})
		}
		for i := 0; i < count; i++ {
			customer := faker.RandomCustomer()
			key := "customer:" + strconv.Itoa(i)
			err := redisutil.StoreJSON(client, key, customer)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		}
		return c.JSON(fiber.Map{"status": "ok", "stored": count})
	}
}

func SearchCustomersHandler(redisURL string) fiber.Handler {
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
		results, err := redisutil.SearchFTSWithLimit(client, "customerIdx", query, limit, offset)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"results": results})
	}
}

func RandomCustomerHandler(redisURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "internal server error: redis client"})
		}
		iter := client.Scan(ctx, 0, "customer:*", 1000).Iterator()
		var keys []string
		for iter.Next(ctx) {
			keys = append(keys, iter.Val())
		}
		if len(keys) == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "no customers found"})
		}
		key := keys[rand.Intn(len(keys))]
		val, err := client.Do(ctx, "JSON.GET", key, "$" ).Text()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		var arr []interface{}
		if err := json.Unmarshal([]byte(val), &arr); err == nil && len(arr) > 0 {
			return c.JSON(arr[0])
		}
		return c.SendString(val)
	}
}

