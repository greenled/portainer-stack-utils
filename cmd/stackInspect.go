package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/greenled/portainer-stack-utils/client"
	portainer "github.com/portainer/portainer/api"

	"github.com/greenled/portainer-stack-utils/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// stackInspectCmd represents the stack inspect command
var stackInspectCmd = &cobra.Command{
	Use:     "inspect",
	Short:   "Inspect a stack",
	Example: "  psu stack inspect mystack",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stackName := args[0]
		var endpointSwarmClusterId string
		var stack portainer.Stack

		var endpoint portainer.Endpoint
		if endpointName := viper.GetString("stack.inspect.endpoint"); endpointName == "" {
			// Guess endpoint if not set
			logrus.WithFields(logrus.Fields{
				"implications": "Command will fail if there is not exactly one endpoint available",
			}).Warning("Endpoint not set")
			var endpointRetrievalErr error
			endpoint, endpointRetrievalErr = common.GetDefaultEndpoint()
			common.CheckError(endpointRetrievalErr)
			endpointName = endpoint.Name
			logrus.WithFields(logrus.Fields{
				"endpoint": endpointName,
			}).Debug("Using the only available endpoint")
		} else {
			// Get endpoint by name
			var endpointRetrievalErr error
			endpoint, endpointRetrievalErr = common.GetEndpointByName(endpointName)
			common.CheckError(endpointRetrievalErr)
		}

		var selectionErr, stackRetrievalErr error
		endpointSwarmClusterId, selectionErr = common.GetEndpointSwarmClusterId(endpoint.ID)
		if selectionErr == nil {
			// It's a swarm cluster
			logrus.WithFields(logrus.Fields{
				"stack":    stackName,
				"endpoint": endpoint.Name,
			}).Debug("Getting stack")
			stack, stackRetrievalErr = common.GetStackByName(stackName, endpointSwarmClusterId, endpoint.ID)
		} else if selectionErr == common.ErrStackClusterNotFound {
			// It's not a swarm cluster
			logrus.WithFields(logrus.Fields{
				"stack":    stackName,
				"endpoint": endpoint.Name,
			}).Debug("Getting stack")
			stack, stackRetrievalErr = common.GetStackByName(stackName, "", endpoint.ID)
		} else {
			// Something else happened
			common.CheckError(selectionErr)
		}

		if stackRetrievalErr == nil {
			// The stack exists
			switch viper.GetString("stack.inspect.format") {
			case "table":
				// Print stack in a table format
				writer, err := common.NewTabWriter([]string{
					"ID",
					"NAME",
					"TYPE",
					"ENDPOINT",
				})
				common.CheckError(err)
				_, err = fmt.Fprintln(writer, fmt.Sprintf(
					"%v\t%s\t%v\t%s",
					stack.ID,
					stack.Name,
					client.GetTranslatedStackType(stack),
					endpoint.Name,
				))
				common.CheckError(err)
				flushErr := writer.Flush()
				common.CheckError(flushErr)
			case "json":
				// Print stack in a json format
				stackJsonBytes, err := json.Marshal(stack)
				common.CheckError(err)
				fmt.Println(string(stackJsonBytes))
			default:
				// Print stack in a custom format
				template, templateParsingErr := template.New("stackTpl").Parse(viper.GetString("stack.inspect.format"))
				common.CheckError(templateParsingErr)
				templateExecutionErr := template.Execute(os.Stdout, stack)
				common.CheckError(templateExecutionErr)
				fmt.Println()
			}
		} else if stackRetrievalErr == common.ErrStackNotFound {
			// The stack does not exist
			logrus.WithFields(logrus.Fields{
				"stack":    stackName,
				"endpoint": endpoint.Name,
			}).Fatal("Stack not found")
		} else {
			// Something else happened
			common.CheckError(stackRetrievalErr)
		}
	},
}

func init() {
	stackCmd.AddCommand(stackInspectCmd)

	stackInspectCmd.Flags().String("endpoint", "", "Filter by endpoint name.")
	stackInspectCmd.Flags().String("format", "table", `Output format. Can be "table", "json" or a Go template.`)
	viper.BindPFlag("stack.inspect.endpoint", stackInspectCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("stack.inspect.format", stackInspectCmd.Flags().Lookup("format"))

	stackInspectCmd.SetUsageTemplate(stackInspectCmd.UsageTemplate() + common.GetFormatHelp(portainer.Stack{}))
}
