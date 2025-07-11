package redisutil

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis URL: %w", err)
	}
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
		"$.primaryIdentifiers.cmec_visitor_id", "AS", "visitor_id", "TEXT",
	).Err()
}

func CreateEventIndex(client *redis.Client) error {
	// Drop index if exists
	client.Do(ctx, "FT.DROPINDEX", "eventIdx")
	// Create index
	return client.Do(ctx, "FT.CREATE", "eventIdx", "ON", "JSON", "PREFIX", "1", "event:",
		"SCHEMA",
		"$.identifiers.cmec_visitor_id", "AS", "visitor_id", "TEXT",
		"$.identifiers.cmec_contact_call_id", "AS", "call_id", "TEXT",
		"$.identifiers.cmec_contact_chat_id", "AS", "chat_id", "TEXT",
		"$.identifiers.cmec_contact_external_id", "AS", "external_id", "TEXT",
		"$.identifiers.cmec_contact_form2lead_id", "AS", "form2lead_id", "TEXT",
		"$.identifiers.cmec_contact_tickets_id", "AS", "tickets_id", "TEXT",
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
