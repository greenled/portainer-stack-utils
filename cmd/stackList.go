package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/greenled/portainer-stack-utils/client"
	portainer "github.com/portainer/portainer/api"

	"github.com/sirupsen/logrus"

	"github.com/greenled/portainer-stack-utils/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// stackListCmd represents the remove command
var stackListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List stacks",
	Aliases: []string{"ls"},
	Example: `  Print stacks in endpoint with name=primary in a table format:
  psu stack ls --endpoint primary

  Print names of stacks in endpoint with name=primary:
  psu stack ls --endpoint primary --format "{{ .Name }}"

  Print environment variables of stacks in all endpoints:
  psu stack ls --format "{{ .Name }}: {{ range .Env }}{{ .Name }}=\"{{ .Value }}\" {{ end }}"`,
	Run: func(cmd *cobra.Command, args []string) {
		portainerClient, err := common.GetClient()
		common.CheckError(err)

		endpoints, endpointsRetrievalErr := portainerClient.EndpointList()
		common.CheckError(endpointsRetrievalErr)

		var endpointSwarmClusterID string
		var stacks []portainer.Stack
		if endpointName := viper.GetString("stack.list.endpoint"); endpointName != "" {
			// Get endpoint by name
			endpoint, endpointRetrievalErr := common.GetEndpointFromListByName(endpoints, endpointName)
			common.CheckError(endpointRetrievalErr)

			logrus.WithFields(logrus.Fields{
				"endpoint": endpoint.Name,
			}).Debug("Getting endpoint's Docker info")
			var selectionErr error
			endpointSwarmClusterID, selectionErr = common.GetEndpointSwarmClusterID(endpoint.ID)
			if selectionErr == nil {
				// It's a swarm cluster
				logrus.WithFields(logrus.Fields{
					"endpoint": endpoint.Name,
				}).Debug("Getting stacks")
				stacks, err = portainerClient.GetStacks(endpointSwarmClusterID, endpoint.ID)
				common.CheckError(err)
			} else if selectionErr == common.ErrStackClusterNotFound {
				// It's not a swarm cluster
				logrus.WithFields(logrus.Fields{
					"endpoint": endpoint.Name,
				}).Debug("Getting stacks")
				stacks, err = portainerClient.GetStacks("", endpoint.ID)
				common.CheckError(err)
			} else {
				// Something else happened
				common.CheckError(selectionErr)
			}
		} else {
			logrus.Debug("Getting stacks")
			stacks, err = portainerClient.GetStacks("", 0)
			common.CheckError(err)
		}

		switch viper.GetString("stack.list.format") {
		case "table":
			// Print stacks in a table format
			writer, err := common.NewTabWriter([]string{
				"ID",
				"NAME",
				"TYPE",
				"ENDPOINT",
			})
			common.CheckError(err)
			for _, s := range stacks {
				stackEndpoint, err := common.GetEndpointFromListByID(endpoints, s.EndpointID)
				common.CheckError(err)
				_, err = fmt.Fprintln(writer, fmt.Sprintf(
					"%v\t%s\t%v\t%s",
					s.ID,
					s.Name,
					client.GetTranslatedStackType(s.Type),
					stackEndpoint.Name,
				))
				common.CheckError(err)
			}
			flushErr := writer.Flush()
			common.CheckError(flushErr)
		case "json":
			// Print stacks in a json format
			stacksJSONBytes, err := json.Marshal(stacks)
			common.CheckError(err)
			fmt.Println(string(stacksJSONBytes))
		default:
			// Print stacks in a custom format
			template, templateParsingErr := template.New("stackTpl").Parse(viper.GetString("stack.list.format"))
			common.CheckError(templateParsingErr)
			for _, s := range stacks {
				templateExecutionErr := template.Execute(os.Stdout, s)
				common.CheckError(templateExecutionErr)
				fmt.Println()
			}
		}
	},
}

func init() {
	stackCmd.AddCommand(stackListCmd)

	stackListCmd.Flags().String("endpoint", "", "Filter by endpoint name.")
	stackListCmd.Flags().String("format", "table", `Output format. Can be "table", "json" or a Go template.`)
	viper.BindPFlag("stack.list.endpoint", stackListCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("stack.list.format", stackListCmd.Flags().Lookup("format"))

	stackListCmd.SetUsageTemplate(stackListCmd.UsageTemplate() + common.GetFormatHelp(portainer.Stack{}))
}
