package cmd

import (
	"github.com/spf13/cobra"
)

// networkCmd represents the network command
var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage networks",
}

func init() {
	rootCmd.AddCommand(networkCmd)
}
