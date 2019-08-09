package cmd

import (
	"github.com/spf13/cobra"
)

// endpointGroupCmd represents the endpoint group command
var endpointGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage endpoint groups",
}

func init() {
	endpointCmd.AddCommand(endpointGroupCmd)
}
