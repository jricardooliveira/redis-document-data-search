package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jricardooliveira/redis-document-data-search/internal/faker"
	"github.com/jricardooliveira/redis-document-data-search/internal/redisutil"
)

// Configuration via environment variables:
//   REDIS_URL   - Redis connection string (default: redis://localhost:6379/0)
//   API_PORT    - HTTP server port (default: 8080)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		dur := time.Since(start)
		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Path()
		if err != nil {
			slog.Error("request error", "method", method, "path", path, "status", status, "duration_μs", dur.Microseconds(), "error", err.Error())
		} else {
			slog.Info("request", "method", method, "path", path, "status", status, "duration_μs", dur.Microseconds())
		}
		return err
	})

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		println("\nShutting down server...")
		_ = app.Shutdown()
		cancel()
	}()

	// Helper to pretty-print any response as HTML <pre>
	prettyJSON := func(c *fiber.Ctx, v interface{}) error {
		pretty, _ := json.MarshalIndent(v, "", "  ")
		return c.Type("html", "utf-8").SendString("<pre>" + string(pretty) + "</pre>")
	}

	app.Post("/generate_customers", func(c *fiber.Ctx) error {
		count, _ := strconv.Atoi(c.Query("count", "1000"))
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			slog.Error("failed to create Redis client", "error", err)
			return prettyJSON(c, fiber.Map{"error": "internal server error: redis client"})
		}
		for i := 0; i < count; i++ {
			customer := faker.RandomCustomer()
			key := "customer:" + strconv.Itoa(i)
			err := redisutil.StoreJSON(client, key, customer)
			if err != nil {
				return prettyJSON(c, fiber.Map{"error": err.Error()})
			}
		}
		return prettyJSON(c, fiber.Map{"status": "ok", "stored": count})
	})

	app.Post("/generate_events", func(c *fiber.Ctx) error {
		count, _ := strconv.Atoi(c.Query("count", "1000"))
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			slog.Error("failed to create Redis client", "error", err)
			return prettyJSON(c.Status(500), fiber.Map{"error": "internal server error: redis client"})
		}
		for i := 0; i < count; i++ {
			event := faker.RandomEvent()
			key := "event:" + strconv.Itoa(i)
			err := redisutil.StoreJSON(client, key, event)
			if err != nil {
				return prettyJSON(c, fiber.Map{"error": err.Error()})
			}
		}
		return prettyJSON(c, fiber.Map{"status": "ok", "stored": count})
	})

	app.Post("/create_indexes", func(c *fiber.Ctx) error {
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			slog.Error("failed to create Redis client", "error", err)
			return prettyJSON(c, fiber.Map{"error": "internal server error: redis client"})
		}
		err1 := redisutil.CreateCustomerIndex(client)
		err2 := redisutil.CreateEventIndex(client)
		if err1 != nil || err2 != nil {
			return prettyJSON(c, fiber.Map{"customerIdx": err1, "eventIdx": err2})
		}
		return prettyJSON(c, fiber.Map{"status": "ok"})
	})

	app.Get("/search_customers", func(c *fiber.Ctx) error {
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			slog.Error("failed to create Redis client", "error", err)
			return prettyJSON(c, fiber.Map{"error": "internal server error: redis client"})
		}
		identifiers := map[string]string{}
		for k, v := range c.Queries() {
			if k != "limit" && k != "offset" {
				identifiers[k] = v
			}
		}
		limit, err1 := strconv.Atoi(c.Query("limit", "10"))
		if err1 != nil || limit < 1 {
			return prettyJSON(c.Status(400), fiber.Map{"error": "limit must be a positive integer"})
		}
		offset, err2 := strconv.Atoi(c.Query("offset", "0"))
		if err2 != nil || offset < 0 {
			return prettyJSON(c.Status(400), fiber.Map{"error": "offset must be a non-negative integer"})
		}
		query := apiBuildRediSearchQuery(identifiers)
		start := time.Now()
		results, err := redisutil.SearchFTSWithLimit(client, "customerIdx", query, limit, offset)
		elapsed := time.Since(start).Microseconds()
		if err != nil {
			return prettyJSON(c, fiber.Map{"error": err.Error(), "query_time_μs": elapsed})
		}
		return prettyJSON(c, fiber.Map{"results": results, "query_time_μs": elapsed})
	})

	app.Get("/search_events", func(c *fiber.Ctx) error {
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			slog.Error("failed to create Redis client", "error", err)
			return prettyJSON(c, fiber.Map{"error": "internal server error: redis client"})
		}
		identifiers := map[string]string{}
		for k, v := range c.Queries() {
			if k != "limit" && k != "offset" {
				identifiers[k] = v
			}
		}
		limit, err1 := strconv.Atoi(c.Query("limit", "10"))
		if err1 != nil || limit < 1 {
			return prettyJSON(c.Status(400), fiber.Map{"error": "limit must be a positive integer"})
		}
		offset, err2 := strconv.Atoi(c.Query("offset", "0"))
		if err2 != nil || offset < 0 {
			return prettyJSON(c.Status(400), fiber.Map{"error": "offset must be a non-negative integer"})
		}
		query := apiBuildRediSearchQuery(identifiers)
		start := time.Now()
		results, err := redisutil.SearchFTSWithLimit(client, "eventIdx", query, limit, offset)
		elapsed := time.Since(start).Microseconds()
		if err != nil {
			return prettyJSON(c, fiber.Map{"error": err.Error(), "query_time_μs": elapsed})
		}
		return prettyJSON(c, fiber.Map{"results": results, "query_time_μs": elapsed})
	})

	app.Get("/random_event", func(c *fiber.Ctx) error {
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			slog.Error("failed to create Redis client", "error", err)
			return prettyJSON(c, fiber.Map{"error": "internal server error: redis client"})
		}
		// Scan for event keys
		iter := client.Scan(ctx, 0, "event:*", 1000).Iterator()
		var keys []string
		for iter.Next(ctx) {
			keys = append(keys, iter.Val())
		}
		if len(keys) == 0 {
			return prettyJSON(c, fiber.Map{"error": "no events found"})
		}
		// Pick a random key
		key := keys[rand.Intn(len(keys))]
		val, err := client.Do(ctx, "JSON.GET", key, "$").Text()
		if err != nil {
			return prettyJSON(c, fiber.Map{"error": err.Error()})
		}
		// Return the JSON (as array, so unmarshal to []interface{} and return the first)
		var arr []interface{}
		if err := json.Unmarshal([]byte(val), &arr); err == nil && len(arr) > 0 {
			return prettyJSON(c, arr[0])
		}
		return c.SendString(val)
	})

	app.Get("/random_customer", func(c *fiber.Ctx) error {
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			slog.Error("failed to create Redis client", "error", err)
			return prettyJSON(c, fiber.Map{"error": "internal server error: redis client"})
		}
		// Scan for customer keys
		iter := client.Scan(ctx, 0, "customer:*", 1000).Iterator()
		var keys []string
		for iter.Next(ctx) {
			keys = append(keys, iter.Val())
		}
		if len(keys) == 0 {
			return prettyJSON(c, fiber.Map{"error": "no customers found"})
		}
		// Pick a random key
		key := keys[rand.Intn(len(keys))]
		val, err := client.Do(ctx, "JSON.GET", key, "$").Text()
		if err != nil {
			return prettyJSON(c, fiber.Map{"error": err.Error()})
		}
		// Return the JSON (as array, so unmarshal to []interface{} and return the first)
		var arr []interface{}
		if err := json.Unmarshal([]byte(val), &arr); err == nil && len(arr) > 0 {
			return prettyJSON(c, arr[0])
		}
		return c.SendString(val)
	})

	app.Get("/healthz", func(c *fiber.Ctx) error {
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			return prettyJSON(c, fiber.Map{"status": "error", "error": "redis client failed", "redis_url": redisURL})
		}
		// Count customer records
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
		// Parse DB index from redisURL
		dbIdx := "0"
		if parts := strings.Split(redisURL, "/"); len(parts) > 1 {
			dbIdx = parts[len(parts)-1]
		}
		return prettyJSON(c, fiber.Map{
			"status":         "ok",
			"redis_url":      redisURL,
			"db_index":       dbIdx,
			"customer_count": customerCount,
			"event_count":    eventCount,
		})
	})

	logger.Info("server starting", "port", port)
	if err := app.Listen(":" + port); err != nil {
		logger.Error("server error", "err", err)
	}
}

func apiBuildRediSearchQuery(identifiers map[string]string) string {
	var parts []string
	for k, v := range identifiers {
		parts = append(parts, "@"+k+":"+v)
	}
	return strings.Join(parts, " ")
}
