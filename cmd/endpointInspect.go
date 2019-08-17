package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	portainer "github.com/portainer/portainer/api"
	"github.com/sirupsen/logrus"

	"github.com/greenled/portainer-stack-utils/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// endpointInspectCmd represents the endpoint inspect command
var endpointInspectCmd = &cobra.Command{
	Use:     "inspect",
	Short:   "Inspect an endpoint",
	Example: "psu endpoint inspect primary",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var endpoint portainer.Endpoint
		if len(args) == 1 {
			var endpointRetrievalErr error
			endpoint, endpointRetrievalErr = common.GetEndpointByName(args[0])
			common.CheckError(endpointRetrievalErr)
		} else {
			logrus.WithFields(logrus.Fields{
				"implications": "Command will fail if there is not exactly one endpoint available",
			}).Warning("Endpoint not set")
			var endpointRetrievalErr error
			endpoint, endpointRetrievalErr = common.GetDefaultEndpoint()
			common.CheckError(endpointRetrievalErr)
			logrus.WithFields(logrus.Fields{
				"endpoint": endpoint.Name,
			}).Debug("Using the only available endpoint")
		}

		switch viper.GetString("endpoint.inspect.format") {
		case "table":
			// Print endpoint in a table format
			writer, err := common.NewTabWriter([]string{
				"ID",
				"NAME",
				"TYPE",
				"URL",
				"PUBLIC URL",
				"GROUP ID",
			})
			common.CheckError(err)
			var endpointType string
			if endpoint.Type == 1 {
				endpointType = "docker"
			} else if endpoint.Type == 2 {
				endpointType = "agent"
			}
			_, err = fmt.Fprintln(writer, fmt.Sprintf(
				"%v\t%s\t%v\t%s\t%s\t%v",
				endpoint.ID,
				endpoint.Name,
				endpointType,
				endpoint.URL,
				endpoint.PublicURL,
				endpoint.GroupID,
			))
			common.CheckError(err)
			err = writer.Flush()
			common.CheckError(err)
		case "json":
			// Print endpoint in a json format
			endpointJsonBytes, err := json.Marshal(endpoint)
			common.CheckError(err)
			fmt.Println(string(endpointJsonBytes))
		default:
			// Print endpoint in a custom format
			template, err := template.New("endpointTpl").Parse(viper.GetString("endpoint.inspect.format"))
			common.CheckError(err)
			err = template.Execute(os.Stdout, endpoint)
			common.CheckError(err)
			fmt.Println()
		}
	},
}

func init() {
	endpointCmd.AddCommand(endpointInspectCmd)

	endpointInspectCmd.Flags().String("format", "table", `Output format. Can be "table", "json" or a Go template.`)
	viper.BindPFlag("endpoint.inspect.format", endpointInspectCmd.Flags().Lookup("format"))

	endpointInspectCmd.SetUsageTemplate(endpointInspectCmd.UsageTemplate() + common.GetFormatHelp(portainer.Endpoint{}))
}
