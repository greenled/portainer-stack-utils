package cmd

import (
	"github.com/spf13/cobra"
)

// Command2 represents the stack command
var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Manage stacks",
}

func init() {
	rootCmd.AddCommand(stackCmd)
}
