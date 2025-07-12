package cli

import (
	"os"

	"github.com/jricardooliveira/redis-document-data-search/internal/cli/commands"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "redisdoccli",
	Short: "Redis Document Data CLI",
}

func init() {
	rootCmd.AddCommand(commands.GenerateCustomersCmd)
	rootCmd.AddCommand(commands.GenerateEventsCmd)
	rootCmd.AddCommand(commands.CreateIndexesCmd)
	rootCmd.AddCommand(commands.SearchCustomersCmd)
	rootCmd.AddCommand(commands.SearchEventsCmd)
	rootCmd.AddCommand(commands.CustomerCmd)
	rootCmd.AddCommand(commands.EventCmd)
	rootCmd.AddCommand(commands.SampleToCSVCommand)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
