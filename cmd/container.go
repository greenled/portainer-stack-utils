package cmd

import (
	"github.com/spf13/cobra"
)

// containerCmd represents the container command
var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Manage containers",
}

func init() {
	rootCmd.AddCommand(containerCmd)
}
