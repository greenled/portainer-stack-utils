package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/greenled/portainer-stack-utils/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check Portainer status",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := common.GetClient()
		common.CheckError(err)

		respBody, err := client.GetStatus()
		common.CheckError(err)

		if viper.GetString("status.format") != "" {
			// Print stack fields formatted
			template, templateParsingErr := template.New("statusTpl").Parse(viper.GetString("status.format"))
			common.CheckError(templateParsingErr)
			templateExecutionErr := template.Execute(os.Stdout, respBody)
			common.CheckError(templateExecutionErr)
			fmt.Println()
		} else {
			// Print status fields as a table
			writer, newTabWriterErr := common.NewTabWriter([]string{
				"VERSION",
				"AUTHENTICATION",
				"ENDPOINT MANAGEMENT",
				"ANALYTICS",
			})
			common.CheckError(newTabWriterErr)

			_, printingErr := fmt.Fprintln(writer, fmt.Sprintf(
				"%s\t%v\t%v\t%v",
				respBody.Version,
				respBody.Authentication,
				respBody.EndpointManagement,
				respBody.Analytics,
			))
			common.CheckError(printingErr)

			flushErr := writer.Flush()
			common.CheckError(flushErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().String("format", "", "Format output using a Go template.")
	viper.BindPFlag("status.format", statusCmd.Flags().Lookup("format"))
}
