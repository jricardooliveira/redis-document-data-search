package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"github.com/spf13/cobra"
	"github.com/jricardooliveira/redis-document-data-search/internal/redisutil"
	"github.com/jricardooliveira/redis-document-data-search/internal/cliutil"
)


var SearchEventsCmd = &cobra.Command{
	Use:   "search_events [key=value ...]",
	Short: "Search events in Redis",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		identifiers := map[string]string{}
		for _, arg := range args {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				identifiers[parts[0]] = parts[1]
			}
		}
		redisURL := os.Getenv("REDIS_URL")
		if redisURL == "" {
			redisURL = "redis://localhost:6379/0"
		}
		client, err := redisutil.NewRedisClient(redisURL)
		if err != nil {
			fmt.Println("Error creating Redis client:", err)
			os.Exit(1)
		}
		query := cliutil.BuildRediSearchQuery(identifiers)
		results, err := redisutil.SearchFTS(client, "eventIdx", query)
		if err != nil {
			fmt.Println("RediSearch error:", err)
			os.Exit(1)
		}
		out, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(out))
	},
}
