package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate Bash completion scripts",
	Example: `  Load completions in the current Bash shell:
  . <(psu completion)

  Configure Bash shell to load completions for each session:
  # ~/.bashrc or ~/.profile
  . <(psu completion)`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
