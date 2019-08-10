package cmd

import (
	"fmt"

	"github.com/greenled/portainer-stack-utils/common"
	portainer "github.com/portainer/portainer/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// stackRemoveCmd represents the remove command
var stackRemoveCmd = &cobra.Command{
	Use:     "remove STACK_NAME",
	Short:   "Remove a stack",
	Aliases: []string{"rm", "down"},
	Example: "psu stack rm mystack",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		portainerClient, clientRetrievalErr := common.GetClient()
		common.CheckError(clientRetrievalErr)

		stackName := args[0]
		endpointId := portainer.EndpointID(viper.GetInt("stack.remove.endpoint"))
		var endpointSwarmClusterId string
		var stack portainer.Stack

		// Guess EndpointID if not set
		if endpointId == 0 {
			logrus.WithFields(logrus.Fields{
				"implications": "Command will fail if there is not exactly one endpoint available",
			}).Warning("Endpoint ID not set")
			endpoint, err := common.GetDefaultEndpoint()
			common.CheckError(err)
			endpointId = endpoint.ID
			logrus.WithFields(logrus.Fields{
				"endpoint": endpointId,
			}).Debug("Using the only available endpoint")
		}

		var selectionErr, stackRetrievalErr error
		endpointSwarmClusterId, selectionErr = common.GetEndpointSwarmClusterId(endpointId)
		switch selectionErr.(type) {
		case nil:
			// It's a swarm cluster
			logrus.WithFields(logrus.Fields{
				"stack":    stackName,
				"endpoint": endpointId,
			}).Debug("Getting stack")
			stack, stackRetrievalErr = common.GetStackByName(stackName, endpointSwarmClusterId, endpointId)
		case *common.StackClusterNotFoundError:
			// It's not a swarm cluster
			logrus.WithFields(logrus.Fields{
				"stack":    stackName,
				"endpoint": endpointId,
			}).Debug("Getting stack")
			stack, stackRetrievalErr = common.GetStackByName(stackName, "", endpointId)
		default:
			// Something else happened
			common.CheckError(selectionErr)
		}

		switch stackRetrievalErr.(type) {
		case nil:
			// The stack exists
			stackId := stack.ID

			logrus.WithFields(logrus.Fields{
				"stack":    stackName,
				"endpoint": endpointId,
			}).Info("Removing stack")
			err := portainerClient.DeleteStack(stackId)
			common.CheckError(err)
			logrus.WithFields(logrus.Fields{
				"stack":    stack.Name,
				"endpoint": stack.EndpointID,
			}).Info("Stack removed")
		case *common.StackNotFoundError:
			// The stack does not exist
			logrus.WithFields(logrus.Fields{
				"stack":    stackName,
				"endpoint": endpointId,
			}).Debug("Stack not found")
			if viper.GetBool("stack.remove.strict") {
				logrus.WithFields(logrus.Fields{
					"stack":       stackName,
					"endpoint":    endpointId,
					"suggestions": fmt.Sprintf("try with a different endpoint: psu stack rm %s --endpoint ENDPOINT_ID", stackName),
				}).Fatal("stack does not exist")
			}
		default:
			// Something else happened
			common.CheckError(stackRetrievalErr)
		}
	},
}

func init() {
	stackCmd.AddCommand(stackRemoveCmd)

	stackRemoveCmd.Flags().Bool("strict", false, "Fail if stack does not exist.")
	stackRemoveCmd.Flags().Int("endpoint", 0, "Endpoint ID.")
	viper.BindPFlag("stack.remove.strict", stackRemoveCmd.Flags().Lookup("strict"))
	viper.BindPFlag("stack.remove.endpoint", stackRemoveCmd.Flags().Lookup("endpoint"))
}
