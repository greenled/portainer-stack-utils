package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"sort"

	"github.com/spf13/cobra"
)

// configListCmd represents the list command
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configs",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		// Get alphabetically ordered list of config keys
		keys := viper.AllKeys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		// List keys and values
		for _, key := range keys {
			fmt.Printf("%s: %v\n", key, viper.Get(key))
		}
	},
}

func init() {
	configCmd.AddCommand(configListCmd)
}
