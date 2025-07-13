// Package redisutil provides utility functions for connecting to Redis, storing and retrieving JSON documents,
// and general Redis operations for the CLI and API in this project.
//
// It abstracts RedisJSON and RediSearch usage, allowing the application to store, query, and manage
// customer and event data as JSON documents in Redis. The functions here are used by both the API handlers
// and CLI commands for all Redis interactions, including creating clients, storing data, and index management.
//
// Index Generation:
//   Indexes are created using RediSearch to enable fast querying and filtering of customer and event JSON documents
//   by key identifiers such as email, phone, visitor_id, call_id, and others. Without these indexes, searching for
//   specific customers or events would require scanning all documents, which is slow and inefficient. Index creation
//   is required for the full-text and field-level search features provided by the API and CLI, allowing for rapid
//   lookups and complex queries over large datasets.
//
// Typical usage:
//   client, err := redisutil.NewRedisClient(redisURL)
//   err := redisutil.StoreJSON(client, "customer:123", customerObj)
//   err := redisutil.CreateCustomerIndex(client)
//   err := redisutil.CreateEventIndex(client)
//
// This package expects Redis to have RedisJSON and RediSearch modules enabled.

package redisutil

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis URL: %w", err)
	}
	opt.PoolSize = 512 // Limita o pool de conexões a 512
	opt.DialTimeout = 5 * time.Second
	opt.ReadTimeout = 3 * time.Second
	opt.WriteTimeout = 3 * time.Second
	opt.PoolTimeout = 4 * time.Second
	return redis.NewClient(opt), nil
}

func StoreJSON(client *redis.Client, key string, value interface{}) error {
	// If value is a string, try to unmarshal to map[string]interface{} to avoid storing stringified JSON
	if str, ok := value.(string); ok {
		var m map[string]interface{}
		err := json.Unmarshal([]byte(str), &m)
		if err == nil {
			slog.Warn("StoreJSON: value was a string, auto-converted to JSON object", "key", key)
			value = m
		} else {
			slog.Error("StoreJSON: value is a string but not valid JSON, storing as string", "key", key)
		}
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// RedisJSON: JSON.SET key $ data
	return client.Do(ctx, "JSON.SET", key, "$", string(data)).Err()
}

func CreateCustomerIndex(client *redis.Client) error {
	// Drop index if exists
	client.Do(ctx, "FT.DROPINDEX", "customerIdx")
	// Create index
	return client.Do(ctx, "FT.CREATE", "customerIdx", "ON", "JSON", "PREFIX", "1", "customer:",
		"SCHEMA",
		"$.primaryIdentifiers.email", "AS", "email", "TEXT",
		"$.primaryIdentifiers.phone", "AS", "phone", "TEXT",
		"$.primaryIdentifiers.visitor_id", "AS", "visitor_id", "TEXT",
	).Err()
}

func CreateEventIndex(client *redis.Client) error {
	// Drop index if exists
	client.Do(ctx, "FT.DROPINDEX", "eventIdx")
	// Create index
	return client.Do(ctx, "FT.CREATE", "eventIdx", "ON", "JSON", "PREFIX", "1", "event:",
		"SCHEMA",
		"$.identifiers.visitor_id", "AS", "visitor_id", "TEXT",
		"$.identifiers.call_id", "AS", "call_id", "TEXT",
		"$.identifiers.chat_id", "AS", "chat_id", "TEXT",
		"$.identifiers.external_id", "AS", "external_id", "TEXT",
		"$.identifiers.lead_id", "AS", "lead_id", "TEXT",
		"$.identifiers.tickets_id", "AS", "tickets_id", "TEXT",
	).Err()
}

func SearchFTS(client *redis.Client, index string, query string) ([]json.RawMessage, error) {
	res, err := client.Do(ctx, "FT.SEARCH", index, query, "RETURN", "1", "$").Result()
	if err != nil {
		return nil, err
	}
	// Parse the new map-based response
	resMap, ok := res.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
	results, ok := resMap["results"].([]interface{})
	if !ok {
		return nil, nil
	}
	var out []json.RawMessage
	for _, r := range results {
		resultMap, ok := r.(map[interface{}]interface{})
		if !ok {
			continue
		}
		extra, ok := resultMap["extra_attributes"].(map[interface{}]interface{})
		if !ok {
			continue
		}
		raw, ok := extra["$"].(string)
		if ok {
			out = append(out, json.RawMessage(raw))
		}
	}
	return out, nil
}

func SearchFTSWithLimit(client *redis.Client, index string, query string, limit int, offset int) ([]json.RawMessage, error) {
	res, err := client.Do(ctx, "FT.SEARCH", index, query, "LIMIT", offset, limit, "RETURN", "1", "$").Result()
	if err != nil {
		return nil, err
	}
	// Parse the new map-based response
	resMap, ok := res.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}
	results, ok := resMap["results"].([]interface{})
	if !ok {
		return nil, nil
	}
	var out []json.RawMessage
	for _, r := range results {
		resultMap, ok := r.(map[interface{}]interface{})
		if !ok {
			continue
		}
		extra, ok := resultMap["extra_attributes"].(map[interface{}]interface{})
		if !ok {
			continue
		}
		raw, ok := extra["$"].(string)
		if ok {
			out = append(out, json.RawMessage(raw))
		}
	}
	return out, nil
}
