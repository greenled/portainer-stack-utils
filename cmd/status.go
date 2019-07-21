package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/greenled/portainer-stack-utils/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"os"
	"text/template"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check Portainer status",
	Run: func(cmd *cobra.Command, args []string) {
		reqUrl, parsingErr := url.Parse(fmt.Sprintf("%s/api/status", viper.GetString("url")))
		common.CheckError(parsingErr)

		req, newRequestErr := http.NewRequest(http.MethodGet, reqUrl.String(), nil)
		common.CheckError(newRequestErr)
		headerErr := common.AddAuthorizationHeader(req)
		common.CheckError(headerErr)
		common.PrintDebugRequest("Get status request", req)

		client := common.NewHttpClient()

		resp, requestExecutionErr := client.Do(req)
		common.CheckError(requestExecutionErr)
		common.PrintDebugResponse("Get status response", resp)

		responseErr := common.CheckResponseForErrors(resp)
		common.CheckError(responseErr)

		var respBody common.Status
		decodingErr := json.NewDecoder(resp.Body).Decode(&respBody)
		common.CheckError(decodingErr)

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

	statusCmd.Flags().String("format", "", "format output using a Go template")
	viper.BindPFlag("status.format", statusCmd.Flags().Lookup("format"))
}
