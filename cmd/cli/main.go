package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jricardooliveira/redis-document-data-search/faker"
	"github.com/jricardooliveira/redis-document-data-search/redisutil"
)

func apiMain() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <command> [options]")
		fmt.Println("Commands: store_customers, store_events, search_customers, search_events, create_indexes, customer, event")
		return
	}
	cmd := os.Args[1]
	count := 1000
	if len(os.Args) > 2 && (cmd == "store_customers" || cmd == "store_events") {
		if n, err := parseInt(os.Args[2]); err == nil {
			count = n
		}
	}
	identifiers := map[string]string{}
	if strings.HasPrefix(cmd, "search_") && len(os.Args) > 2 {
		for _, arg := range os.Args[2:] {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				identifiers[parts[0]] = parts[1]
			}
		}
	}

	// Get Redis URL from environment variable or use default
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}

	switch cmd {
	case "store_customers":
		client := redisutil.NewRedisClient(redisURL)
		for i := 0; i < count; i++ {
			customer := faker.RandomCustomer()
			key := fmt.Sprintf("customer:%d", i)
			err := redisutil.StoreJSON(client, key, customer)
			if err != nil {
				fmt.Printf("Error storing customer %d: %v\n", i, err)
			}
			if (i+1)%100 == 0 {
				fmt.Printf("Stored %d customers...\n", i+1)
			}
		}
		fmt.Printf("Done. Stored %d customers in Redis.\n", count)
	case "store_events":
		client := redisutil.NewRedisClient(redisURL)
		for i := 0; i < count; i++ {
			event := faker.RandomEvent()
			key := fmt.Sprintf("event:%d", i)
			err := redisutil.StoreJSON(client, key, event)
			if err != nil {
				fmt.Printf("Error storing event %d: %v\n", i, err)
			}
			if (i+1)%100 == 0 {
				fmt.Printf("Stored %d events...\n", i+1)
			}
		}
		fmt.Printf("Done. Stored %d events in Redis.\n", count)
	case "create_indexes":
		client := redisutil.NewRedisClient(redisURL)
		err1 := redisutil.CreateCustomerIndex(client)
		err2 := redisutil.CreateEventIndex(client)
		if err1 != nil {
			fmt.Println("Error creating customer index:", err1)
		} else {
			fmt.Println("Created customerIdx on email, phone, visitor_id")
		}
		if err2 != nil {
			fmt.Println("Error creating event index:", err2)
		} else {
			fmt.Println("Created eventIdx on all identifiers")
		}
	case "search_customers":
		client := redisutil.NewRedisClient(redisURL)
		query := buildRediSearchQuery(identifiers)
		results, err := redisutil.SearchFTS(client, "customerIdx", query)
		if err != nil {
			fmt.Println("RediSearch error:", err)
			return
		}
		out, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(out))
	case "search_events":
		client := redisutil.NewRedisClient(redisURL)
		query := buildRediSearchQuery(identifiers)
		results, err := redisutil.SearchFTS(client, "eventIdx", query)
		if err != nil {
			fmt.Println("RediSearch error:", err)
			return
		}
		out, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(out))
	case "customer":
		out, _ := json.MarshalIndent(faker.RandomCustomer(), "", "  ")
		fmt.Println(string(out))
	case "event":
		out, _ := json.MarshalIndent(faker.RandomEvent(), "", "  ")
		fmt.Println(string(out))
	default:
		fmt.Println("Unknown command:", cmd)
	}
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}

func buildRediSearchQuery(identifiers map[string]string) string {
	var parts []string
	for k, v := range identifiers {
		parts = append(parts, fmt.Sprintf("@%s:%s", k, v))
	}
	return strings.Join(parts, " ")
}
