package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jricardooliveira/redis-document-data-search/internal/api/handlers"
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
	cancel := func() {}
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

	// Route registrations using refactored handlers
	app.Post("/generate_customers", handlers.GenerateCustomersHandler(redisURL))
	app.Post("/generate_events", handlers.GenerateEventsHandler(redisURL))
	app.Post("/create_indexes", handlers.CreateIndexesHandler(redisURL))
	app.Get("/search_customers", handlers.SearchCustomersHandler(redisURL))
	app.Get("/search_events", handlers.SearchEventsHandler(redisURL))
	app.Get("/random_event", handlers.RandomEventHandler(redisURL))
	app.Get("/random_customer", handlers.RandomCustomerHandler(redisURL))
	app.Get("/healthz", handlers.HealthHandler(redisURL))

	logger.Info("server starting", "port", port)
	if err := app.Listen(":" + port); err != nil {
		logger.Error("server error", "err", err)
	}
}
