package commands

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

// SampleToCSVCommand implements fast sampling of Redis documents to CSV
var SampleToCSVCommand = &cobra.Command{
	Use:   "sample_to_csv",
	Short: "Extract a random sample of customer or event documents from Redis and write to a CSV file",
	RunE: func(cmd *cobra.Command, args []string) error {
		typeStr, _ := cmd.Flags().GetString("type")
		percent, _ := cmd.Flags().GetInt("percent")
		output, _ := cmd.Flags().GetString("output")
		redisURL, _ := cmd.Flags().GetString("redis")
		if typeStr != "customer" && typeStr != "event" {
			return fmt.Errorf("--type must be 'customer' or 'event'")
		}
		if percent <= 0 || percent > 100 {
			return fmt.Errorf("--percent must be between 1 and 100")
		}
		if output == "" {
			return fmt.Errorf("--output is required")
		}
		if redisURL == "" {
			redisURL = "redis://localhost:6379"
		}

		ctx := context.Background()
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			return fmt.Errorf("invalid redis url: %v", err)
		}
		client := redis.NewClient(opt)
		defer client.Close()

		// 1. List all keys
		pattern := fmt.Sprintf("%s:*", typeStr)
		var cursor uint64
		var allKeys []string
		for {
			keys, next, err := client.Scan(ctx, cursor, pattern, 1000).Result()
			if err != nil {
				return fmt.Errorf("failed to scan keys: %v", err)
			}
			allKeys = append(allKeys, keys...)
			if next == 0 {
				break
			}
			cursor = next
		}
		if len(allKeys) == 0 {
			return fmt.Errorf("no keys found for pattern %s", pattern)
		}

		// 2. Randomly sample keys
		sampleSize := len(allKeys) * percent / 100
		if sampleSize < 1 {
			sampleSize = 1
		}
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(allKeys), func(i, j int) { allKeys[i], allKeys[j] = allKeys[j], allKeys[i] })
		sampledKeys := allKeys[:sampleSize]
		sort.Strings(sampledKeys)

		// 3. Fetch documents in batches
		pipe := client.Pipeline()
		cmds := make([]*redis.Cmd, len(sampledKeys))
		for i, key := range sampledKeys {
			cmds[i] = pipe.Do(ctx, "JSON.GET", key, "$")
		}
		_, err = pipe.Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch documents: %v", err)
		}

		// 4. Write CSV
		f, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("failed to create output CSV: %v", err)
		}
		defer f.Close()
		w := csv.NewWriter(f)
		defer w.Flush()

		if typeStr == "customer" {
			header := []string{"key", "email", "phone", "visitor_id"}
			w.Write(header)
		} else {
			header := []string{"key", "visitor_id", "call_id", "chat_id", "external_id", "form2lead_id", "tickets_id"}
			w.Write(header)
		}

		for i, cmd := range cmds {
			jsonStr, err := cmd.Text()
			if err != nil || jsonStr == "" {
				continue // skip missing
			}
			// Remove outer array if present
			jsonStr = strings.TrimSpace(jsonStr)
			if strings.HasPrefix(jsonStr, "[") && strings.HasSuffix(jsonStr, "]") {
				jsonStr = strings.TrimPrefix(jsonStr, "[")
				jsonStr = strings.TrimSuffix(jsonStr, "]")
			}
			var record map[string]interface{}
			if err := jsonUnmarshalCompat(jsonStr, &record); err != nil {
				continue
			}
			row := []string{sampledKeys[i]}
			if typeStr == "customer" {
				row = append(row,
					getStringField(record, "primaryIdentifiers", "email"),
					getStringField(record, "primaryIdentifiers", "phone"),
					getStringField(record, "primaryIdentifiers", "cmec_visitor_id"),
				)
			} else {
				row = append(row,
					getStringField(record, "identifiers", "cmec_visitor_id"),
					getStringField(record, "identifiers", "cmec_contact_call_id"),
					getStringField(record, "identifiers", "cmec_contact_chat_id"),
					getStringField(record, "identifiers", "cmec_contact_external_id"),
					getStringField(record, "identifiers", "cmec_contact_form2lead_id"),
					getStringField(record, "identifiers", "cmec_contact_tickets_id"),
				)
			}
			w.Write(row)
		}

		fmt.Printf("Sampled %d records out of %d. Output written to %s\n", sampleSize, len(allKeys), output)
		return nil
	},
}

func init() {
	SampleToCSVCommand.Flags().String("type", "", "Type of document: customer or event (required)")
	SampleToCSVCommand.Flags().Int("percent", 5, "Percent of records to sample (default 5)")
	SampleToCSVCommand.Flags().String("output", "", "Output CSV file (required)")
	SampleToCSVCommand.Flags().String("redis", "", "Redis connection URL (optional)")

	_ = SampleToCSVCommand.MarkFlagRequired("type")
	_ = SampleToCSVCommand.MarkFlagRequired("output")
}

// Helper to unmarshal JSON using encoding/json or a fallback for compatibility
func jsonUnmarshalCompat(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// Helper to extract nested string fields safely
func getStringField(m map[string]interface{}, parent, field string) string {
	if m == nil {
		return ""
	}
	p, ok := m[parent].(map[string]interface{})
	if !ok {
		return ""
	}
	val, ok := p[field]
	if !ok {
		return ""
	}
	s, _ := val.(string)
	return s
}
