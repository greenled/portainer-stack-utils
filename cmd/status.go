package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/greenled/portainer-stack-utils/util"

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
		util.CheckError(err)

		respBody, err := client.GetStatus()
		util.CheckError(err)

		if viper.GetString("status.format") != "" {
			// Print stack fields formatted
			template, templateParsingErr := template.New("statusTpl").Parse(viper.GetString("status.format"))
			util.CheckError(templateParsingErr)
			templateExecutionErr := template.Execute(os.Stdout, respBody)
			util.CheckError(templateExecutionErr)
			fmt.Println()
		} else {
			// Print status fields as a table
			writer, newTabWriterErr := util.NewTabWriter([]string{
				"VERSION",
				"AUTHENTICATION",
				"ENDPOINT MANAGEMENT",
				"ANALYTICS",
			})
			util.CheckError(newTabWriterErr)

			_, printingErr := fmt.Fprintln(writer, fmt.Sprintf(
				"%s\t%v\t%v\t%v",
				respBody.Version,
				respBody.Authentication,
				respBody.EndpointManagement,
				respBody.Analytics,
			))
			util.CheckError(printingErr)

			flushErr := writer.Flush()
			util.CheckError(flushErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().String("format", "", "format output using a Go template")
	viper.BindPFlag("status.format", statusCmd.Flags().Lookup("format"))
}
