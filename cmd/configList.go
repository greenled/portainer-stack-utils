package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// configListCmd represents the list command
var configListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List configs",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		// Get alphabetically ordered list of config keys
		keys := viper.AllKeys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		for _, key := range keys {
			// List config key and value
			fmt.Printf("%s: %v\n", key, viper.Get(key))
		}
	},
}

func init() {
	configCmd.AddCommand(configListCmd)
}
