package commands

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/jricardooliveira/redis-document-data-search/internal/faker"
)

var CustomerCmd = &cobra.Command{
	Use:   "customer",
	Short: "Print a random customer JSON",
	Run: func(cmd *cobra.Command, args []string) {
		out, _ := json.MarshalIndent(faker.RandomCustomer(), "", "  ")
		fmt.Println(string(out))
	},
}
