package cmd

import (
	"github.com/spf13/cobra"
)

// secretCmd represents the secret command
var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage secrets",
}

func init() {
	rootCmd.AddCommand(secretCmd)
}
