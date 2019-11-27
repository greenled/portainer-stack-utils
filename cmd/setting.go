package cmd

import (
	"github.com/spf13/cobra"
)

// settingCmd represents the config command
var settingCmd = &cobra.Command{
	Use:   "setting",
	Short: "Manage settings",
}

func init() {
	rootCmd.AddCommand(settingCmd)
}
