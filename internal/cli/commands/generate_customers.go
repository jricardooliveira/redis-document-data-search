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

var GenerateCustomersCmd = &cobra.Command{
	Use:   "generate_customers [count]",
	Short: "Generate and store random customers in Redis",
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
			customer := faker.RandomCustomer()
			key := "customer:" + strconv.Itoa(i)
			err := redisutil.StoreJSON(client, key, customer)
			if err != nil {
				fmt.Printf("Error storing customer %d: %v\n", i, err)
			}
			if (i+1)%100 == 0 {
				fmt.Printf("Stored %d customers...\n", i+1)
			}
		}
		fmt.Printf("Done. Stored %d customers in Redis.\n", count)
	},
}
