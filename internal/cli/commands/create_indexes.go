package commands

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/jricardooliveira/redis-document-data-search/internal/redisutil"
)

var CreateIndexesCmd = &cobra.Command{
	Use:   "create_indexes",
	Short: "Create RediSearch indexes in Redis",
	Run: func(cmd *cobra.Command, args []string) {
		redisURL := os.Getenv("REDIS_URL")
		if redisURL == "" {
			redisURL = "redis://localhost:6379/0"
		}
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			fmt.Println("Error creating Redis client:", err)
			os.Exit(1)
		}
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
	},
}
