package cmd

import (
	"github.com/spf13/cobra"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage services",
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
