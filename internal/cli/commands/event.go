package commands

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/jricardooliveira/redis-document-data-search/internal/faker"
)

var EventCmd = &cobra.Command{
	Use:   "event",
	Short: "Print a random event JSON",
	Run: func(cmd *cobra.Command, args []string) {
		out, _ := json.MarshalIndent(faker.RandomEvent(), "", "  ")
		fmt.Println(string(out))
	},
}
