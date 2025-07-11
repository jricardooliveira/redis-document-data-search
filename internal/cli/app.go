package cli

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/jricardooliveira/redis-document-data-search/internal/cli/commands"
)

var rootCmd = &cobra.Command{
	Use:   "redisdoccli",
	Short: "Redis Document Data CLI",
}

func init() {
	rootCmd.AddCommand(commands.StoreCustomersCmd)
	rootCmd.AddCommand(commands.StoreEventsCmd)
	rootCmd.AddCommand(commands.CreateIndexesCmd)
	rootCmd.AddCommand(commands.SearchCustomersCmd)
	rootCmd.AddCommand(commands.SearchEventsCmd)
	rootCmd.AddCommand(commands.CustomerCmd)
	rootCmd.AddCommand(commands.EventCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
