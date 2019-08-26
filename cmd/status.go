package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	portainer "github.com/portainer/portainer/api"

	"github.com/greenled/portainer-stack-utils/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check Portainer status",
	Example: `  Print status in a table format:
  psu status

  Print version of Portainer server:
  psu status --format "{{ .Version }}"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := common.GetClient()
		common.CheckError(err)

		respBody, err := client.Status()
		common.CheckError(err)

		switch viper.GetString("status.format") {
		case "table":
			// Print status in a table format
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
		case "json":
			// Print status in a json format
			statusJSONBytes, err := json.Marshal(respBody)
			common.CheckError(err)
			fmt.Println(string(statusJSONBytes))
		default:
			// Print status in a custom format
			template, templateParsingErr := template.New("statusTpl").Parse(viper.GetString("status.format"))
			common.CheckError(templateParsingErr)
			templateExecutionErr := template.Execute(os.Stdout, respBody)
			common.CheckError(templateExecutionErr)
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().String("format", "table", `Output format. Can be "table", "json" or a Go template.`)
	viper.BindPFlag("status.format", statusCmd.Flags().Lookup("format"))

	statusCmd.SetUsageTemplate(statusCmd.UsageTemplate() + common.GetFormatHelp(portainer.Status{}))
}
