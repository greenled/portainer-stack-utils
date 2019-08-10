package cmd

import (
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
	Example: "psu stack list --endpoint 1",
	Run: func(cmd *cobra.Command, args []string) {
		portainerClient, err := common.GetClient()
		common.CheckError(err)

		endpointId := portainer.EndpointID(viper.GetInt("stack.list.endpoint"))
		var endpointSwarmClusterId string
		var stacks []portainer.Stack
		if endpointId != 0 {
			var selectionErr error
			endpointSwarmClusterId, selectionErr = common.GetEndpointSwarmClusterId(endpointId)
			switch selectionErr.(type) {
			case nil:
				// It's a swarm cluster
				logrus.WithFields(logrus.Fields{
					"endpoint": endpointId,
				}).Debug("Getting stacks")
				stacks, err = portainerClient.GetStacks(endpointSwarmClusterId, endpointId)
				common.CheckError(err)
			case *common.StackClusterNotFoundError:
				// It's not a swarm cluster
				logrus.WithFields(logrus.Fields{
					"endpoint": endpointId,
				}).Debug("Getting stacks")
				stacks, err = portainerClient.GetStacks("", endpointId)
				common.CheckError(err)
			default:
				// Something else happened
				common.CheckError(selectionErr)
			}
		} else {
			logrus.Debug("Getting stacks")
			stacks, err = portainerClient.GetStacks("", 0)
			common.CheckError(err)
		}

		if viper.GetString("stack.list.format") != "" {
			// Print stack fields formatted
			template, templateParsingErr := template.New("stackTpl").Parse(viper.GetString("stack.list.format"))
			common.CheckError(templateParsingErr)
			for _, s := range stacks {
				templateExecutionErr := template.Execute(os.Stdout, s)
				common.CheckError(templateExecutionErr)
				fmt.Println()
			}
		} else {
			// Print all stack fields as a table
			writer, err := common.NewTabWriter([]string{
				"ID",
				"NAME",
				"TYPE",
				"ENDPOINT ID",
			})
			common.CheckError(err)
			for _, s := range stacks {
				_, err := fmt.Fprintln(writer, fmt.Sprintf(
					"%v\t%s\t%v\t%v",
					s.ID,
					s.Name,
					client.GetTranslatedStackType(s),
					s.EndpointID,
				))
				common.CheckError(err)
			}
			flushErr := writer.Flush()
			common.CheckError(flushErr)
		}
	},
}

func init() {
	stackCmd.AddCommand(stackListCmd)

	stackListCmd.Flags().Int("endpoint", 0, "Filter by endpoint ID.")
	stackListCmd.Flags().String("format", "", "Format output using a Go template.")
	viper.BindPFlag("stack.list.endpoint", stackListCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("stack.list.format", stackListCmd.Flags().Lookup("format"))
}
