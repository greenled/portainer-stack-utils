package cmd

import (
	"github.com/spf13/cobra"
)

// volumeCmd represents the volume command
var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Manage volumes",
}

func init() {
	rootCmd.AddCommand(volumeCmd)
}
