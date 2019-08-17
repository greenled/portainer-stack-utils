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
	Use:     "remove <name>",
	Short:   "Remove a stack",
	Aliases: []string{"rm", "down"},
	Example: "  psu stack rm mystack",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		portainerClient, clientRetrievalErr := common.GetClient()
		common.CheckError(clientRetrievalErr)

		stackName := args[0]
		var endpointSwarmClusterId string
		var stack portainer.Stack

		var endpoint portainer.Endpoint
		if endpointName := viper.GetString("stack.remove.endpoint"); endpointName == "" {
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
			stackId := stack.ID

			logrus.WithFields(logrus.Fields{
				"stack":    stackName,
				"endpoint": endpoint.Name,
			}).Info("Removing stack")
			err := portainerClient.DeleteStack(stackId)
			common.CheckError(err)
			logrus.WithFields(logrus.Fields{
				"stack":    stack.Name,
				"endpoint": endpoint.Name,
			}).Info("Stack removed")
		} else if stackRetrievalErr == common.ErrStackNotFound {
			// The stack does not exist
			logrus.WithFields(logrus.Fields{
				"stack":    stackName,
				"endpoint": endpoint.Name,
			}).Debug("Stack not found")
			if viper.GetBool("stack.remove.strict") {
				logrus.WithFields(logrus.Fields{
					"stack":       stackName,
					"endpoint":    endpoint.Name,
					"suggestions": fmt.Sprintf("try with a different endpoint: psu stack rm %s --endpoint ENDPOINT_NAME", stackName),
				}).Fatal("stack does not exist")
			}
		} else {
			// Something else happened
			common.CheckError(stackRetrievalErr)
		}
	},
}

func init() {
	stackCmd.AddCommand(stackRemoveCmd)

	stackRemoveCmd.Flags().Bool("strict", false, "Fail if stack does not exist.")
	stackRemoveCmd.Flags().String("endpoint", "", "Endpoint name.")
	viper.BindPFlag("stack.remove.strict", stackRemoveCmd.Flags().Lookup("strict"))
	viper.BindPFlag("stack.remove.endpoint", stackRemoveCmd.Flags().Lookup("endpoint"))
}
