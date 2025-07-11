package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/jricardooliveira/redis-document-data-search/internal/cliutil"
	"github.com/jricardooliveira/redis-document-data-search/internal/faker"
	"github.com/jricardooliveira/redis-document-data-search/internal/redisutil"
)

var GenerateEventsCmd = &cobra.Command{
	Use:   "generate_events [count]",
	Short: "Generate and store random events in Redis",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		count := 1000
		if len(args) > 0 {
			if n, err := cliutil.ParseInt(args[0]); err == nil {
				count = n
			}
		}
		redisURL := os.Getenv("REDIS_URL")
		if redisURL == "" {
			redisURL = "redis://localhost:6379/0"
		}
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			fmt.Println("Error creating Redis client:", err)
			return
		}
		for i := 0; i < count; i++ {
			event := faker.RandomEvent()
			key := "event:" + strconv.Itoa(i)
			err := redisutil.StoreJSON(client, key, event)
			if err != nil {
				fmt.Printf("Error storing event %d: %v\n", i, err)
			}
			if (i+1)%100 == 0 {
				fmt.Printf("Stored %d events...\n", i+1)
			}
		}
		fmt.Printf("Done. Stored %d events in Redis.\n", count)
	},
}
